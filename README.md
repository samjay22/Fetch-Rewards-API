Receipt Processor Project
Overview

The Receipt Processor project is a RESTful API designed to handle receipt submission and retrieval of points awarded based on processed receipts. This document provides an overview of the project structure, API endpoints, data schemas, technologies used, security measures, and deployment details.
OpenAPI Specification

The project adheres to the OpenAPI 3.0.3 standard. Below is a summary of the API endpoints defined:

Endpoints

    Submit Receipt
        Endpoint: POST /receipts/process
        Description: Submits a receipt for processing.
        Request Body: JSON format following the Receipt schema.
        Responses:
            200 OK: Returns the ID assigned to the receipt.
            400 Bad Request: Indicates an invalid receipt.

    Get Points
        Endpoint: GET /receipts/{id}/points
        Description: Retrieves the points awarded for a specific receipt ID.
        Path Parameter:
            id: The ID of the receipt.
        Responses:
            200 OK: Returns the number of points awarded.
            404 Not Found: Indicates no receipt found for the provided ID.

Components

The OpenAPI schema defines the following components:

    Schemas: Receipt and Item describe the structure of receipts and individual items on a receipt, respectively.

Technology Stack

Backend

    Language: Golang
    Design Patterns: Delegate Pattern, SOLID Principles, Event-Driven Architecture, Microservice/Controller Pattern.
    Database: SQLite with SQL for data persistence. The database was created to be clusterd and Indexes were created to speed up data retreval.

Frontend

    Language: Vanilla JavaScript, HTML, and CSS
    Purpose: Provides an interface for users to search and add receipts.

Security Measures

    SSL Certificate: Ensures secure communication between clients and the server.
    Web Proxy: Implemented in C# .NET, hosted on a VPS behind a Cloudflare-managed domain.
    Tunnel Setup: Uses Cloudflare for DNS management and VPN tunnels via No-Trust to secure API requests.

Deployment

    Server: Custom-built server hosted behind my proxy and VPN tunnel into my network.
    Cloudflare: Acts as a reverse proxy, providing an additional layer of security and performance optimization.

Usage

    1: Clone the repo using git clone https://github.com/samjay22/Fetch-Rewards-API/
    2: run cd Fetch-Rewards-API to navigate into the repo directory
    3: Ensure golang version 22 or higher is installed. Run go mod tidy in the directory to obtain dependencies and refresh the project index
    4: Build the go binary by running go build -o MY_EXE_NAME
    
Conclusion

  For other details regarding the project as to technical decisions, questions, or concerns, do not hesitate to reach out via email or request to become a contributor. My email is samjaytaylor22@gmail.com. 
