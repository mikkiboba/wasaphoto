package api

import (
	"errors"
	"io"
	"net/http"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

/*
Handler of the operation PUT for the path /users/:username/posts
It needs the username of the user posting the photo in the parameters and the file posted in the request body.
*/
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	rt.baseLogger.Info("Uploading...")

	// Start the transaction
	StartTransaction(rt, w)

	// Generates a new uuid to name the post -> pid (post id)
	pid, err := uuid.NewV4()
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the generation of the token for the post")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Get the username from the parameters and obtains the user's token
	username := ps.ByName("username")
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Insert the photo into the database
	err = rt.db.UploadPhoto(uid, pid.String())
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("The element was not added to the database")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the photo upload into the database")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Create the post path folder, if it exists already it does nothing
	err = os.MkdirAll(rt.postPath, 0777)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with creating the posts path")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Check if the size of the input file is MAX 32 mb
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the input file")
		w.WriteHeader(http.StatusBadRequest)
		Rollback(rt, w)
		return
	}

	// Generate the photo file
	file, _, err := r.FormFile("file")
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the form of the file")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}
	defer file.Close()

	// Create the file which is going to be saved in the directory
	imgFile, err := os.Create(rt.postPath + pid.String())
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the creation of the photo file in the disk")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}
	defer imgFile.Close()

	// Copy the file into the imgFile
	_, err = io.Copy(imgFile, file)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the copy of the file")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)

		// Since there is an error, delete the imgFile from the disk
		err = os.Remove(rt.postPath + pid.String())
		if err != nil {
			rt.baseLogger.WithError(err).Error("There is an error with the removal of the file from the disk")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Commit of the transaction
	Commit(rt, w)

	rt.baseLogger.Info("The photo has been uploaded with id")

}

/*
Handler of the DELETE operation for the route /users/:username/posts/:postid
*/
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username
	username := ps.ByName("username")

	// Get the id
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Get the post id
	pid := ps.ByName("postid")

	// Check if the user is the owner of the photo
	oid, err := rt.db.GetPostOwner(pid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the post owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if uid != oid {
		rt.baseLogger.Error("A user can't delete a post they don't own")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	StartTransaction(rt, w)

	filename, err := rt.db.DeletePhoto(pid)
	if errors.Is(err, database.ErrElementNotDeleted) {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the post from the database. Maybe the post doesn't exist")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the post from the database")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Remove the file from the disk
	err = os.Remove(rt.postPath + filename)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while deleting the photo from the disk")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
	}

	Commit(rt, w)

	w.WriteHeader(http.StatusNoContent)
	rt.baseLogger.Info("Photo deleted succesfully")
}

/*
Handler of the GET operation for the route
*/
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	pid := ps.ByName("postid")

	filename, err := rt.db.GetPhoto(pid)
	if errors.Is(err, database.ErrPhotoNotFound) {
		rt.baseLogger.WithError(err).Error("Photo not found")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, err := os.Open(rt.postPath + filename)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while opening the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/*")
	_, err = io.Copy(w, file)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error sending the image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Info("Photo sent succesfully")
}
