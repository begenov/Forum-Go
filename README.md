# Web Forum Project

This project is a web forum that allows communication between registered users, association of categories to posts, liking and disliking posts and comments, and filtering posts. The project uses SQLite for data storage, Go programming language, and Docker for containerization.

## Features

* User Authentication: users can register on the forum, login, and create a session that expires after a set duration. Users can't register with an email that's already taken, and passwords are encrypted when stored.

* Communication: registered users can create posts and comments, and associate one or more categories with each post. All users can view posts and comments, but only registered users can create them.

* Likes and Dislikes: only registered users can like or dislike posts and comments. The number of likes and dislikes is visible to all users.

* Filtering: users can filter posts by categories. Filtering by categories is equivalent to subforums.


## Dependencies

   * Go (standard packages)
   * SQLite
   * bcrypt
   * UUID
   * Docker

## Installation

1. Clone the repository.
2.    In the root directory of the project, build the Docker image by running the command docker build -t web-forum ..
3.    After the build is complete, start the container using the command docker run -p 8080:8080 web-forum.

or 

```go run cmd/main.go```


The web forum should now be accessible by visiting localhost:8080 in your web browser.

