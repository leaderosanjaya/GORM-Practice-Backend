# GORM REST API Remote Config Practice
A simple project made to demonstrate REST api calls using GORM

# API Routes
| Scope | Path                                | Method | Function                        | Description                           |
|-------|-------------------------------------|--------|---------------------------------|---------------------------------------|
|       | /api                                | GET    | index                           | Checks if the API is alive or not     |
| Auth  | /api/login                          | POST   | authHandler.login               | Login, receive auth token             |
| User  | /api/users                          | POST   | userHandler.CreateUserHandler   | Create new user                       |
| User  | /auth/api/users/{user_id}           | GET    | userHandler.GetUserByID         | Get User Info                         |
| User  | /auth/api/users/{user_id}           | DELETE | userHandler.DeleteUserHandler   | Delete User                           |
| User  | /auth/api/users/{user_id}/keys      | GET    | keyHandler.GetKeysByUserID      | Get Key by User ID                    |
| Tribe | /auth/api/tribes                    | POST   | tribeHandler.CreateTribeHandler | Create Tribe                          |
| Tribe | /auth/api/tribes/{tribe_id}         | DELETE | tribeHandler.DeleteTribeHandler | Delete Tribe                          |
| Tribe | /auth/api/tribes/{tribe_id}/members | POST   | tribeHandler.AssignUser         | Add members to the tribe              |
| Tribe | /auth/api/tribes/{tribe_id}/members | DELETE | tribeHandler.RemoveAssign       | Remove member to the tribe            |
| Tribe | /auth/api/tribes/{tribe_id}         | GET    | tribeHandler.GetTribeByID       | Get tribe info                        |
| Tribe | /auth/api/tribes/{tribe_id}/keys    | GET    | keyHandler.GetKeysByTribeID     | Get tribe keys info                   |
| Key   | /auth/api/keys                      | POST   | keyHandler.CreateKeyHandler     | Create Key                            |
| Key   | /auth/api/keys                      | GET    | keyHandler.GetKeysHandler       | Get All keys                          |
| Key   | /auth/api/keys/{key_id}             | GET    | keyHandler.GetKeyByID           | Get Key by ID                         |
| Key   | /auth/api/keys/{key_id}             | DELETE | keyHandler.DeleteKeyHandler     | Delete Key                            |
| Key   | /auth/api/keys/{key_id}             | PUT    | keyHandler.UpdateKeyByID        | Update Key                            |
| Key   | /auth/api/keys/{key_id}/shares      | POST   | keyHandler.ShareKey             | Share key to another member           |
| Key   | /auth/api/keys/{key_id}/shares      | DELETE | keyHandler.RevokeShare          | Remove shared key from another member |

#TBA