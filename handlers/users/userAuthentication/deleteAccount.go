package userAuthenticate

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )
// // Simulated user data (replace this with your actual database)
// var users = map[string]User{
// 	"user1": {ID: "user1", Name: "John Doe", Email: "john@example.com"},
// 	"user2": {ID: "user2", Name: "Jane Doe", Email: "jane@example.com"},
// }

// func main() {
// 	router := mux.NewRouter()

// 	router.HandleFunc("/api/delete_account/{id}", DeleteAccount).Methods("DELETE")

// 	log.Fatal(http.ListenAndServe(":8080", router))
// }

// // DeleteAccount handles the deletion of a user account
// func DeleteAccount(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	userID := params["id"]

// 	// Check if the user exists
// 	if _, ok := users[userID]; !ok {
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
// 		return
// 	}

// 	// Perform account deletion (replace this with your actual deletion logic)
// 	delete(users, userID)

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
// }
