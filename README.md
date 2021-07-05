# image-converter
Service that expose a RESTful API to convert JPEG to PNG and vice versa and compress the image 
with the compression ratio specified by the user. The user has the ability to view
the history and status of their requests (for example, queued, processed, completed) and upload 
the original image and the processed one.
# Architecture Diagram
![alt text](./docs/architecture-diagram.png)
# Endpoints
- /requests - view user requests history [GET]
- /user/login - user authorization [POST]
- /user/signup - user registration [POST]
- /user/logout - log out of the user [GET]
- /conversion - convert needed image [POST]
- /images/{id} - get needed image [GET]
# Database Scheme
![alt text](./docs/db.png)
