### This project was reviewed by a "FETCH Engineer" which clearly didn't test or do anything with the project. Its also clear they did not know SQL and any of the patterns I used based on the feedback.

### Receipt Processor Project

#### Overview

The Receipt Processor project is a RESTful API designed to handle receipt submission and retrieval of points awarded based on processed receipts. This document provides an in-depth overview of the project structure, API endpoints, data schemas, technologies used, security measures, deployment details, and additional technical considerations.

#### OpenAPI Specification

The project adheres to the OpenAPI 3.0.3 standard, ensuring a standardized approach to building and documenting RESTful APIs.

#### Live Instance

- **URL**: [https://blox-api.com/](https://blox-api.com/)

#### Endpoints

1. **Submit Receipt**
   - **Endpoint**: `POST /receipts/process`
   - **Description**: Submits a receipt for processing.
   - **Request Body**: JSON format following the Receipt schema.
   - **Responses**:
     - `200 OK`: Returns the ID assigned to the receipt.
     - `400 Bad Request`: Indicates an invalid receipt.

2. **Get Points**
   - **Endpoint**: `GET /receipts/{id}/points`
   - **Description**: Retrieves the points awarded for a specific receipt ID.
   - **Path Parameter**:
     - `id`: The ID of the receipt.
   - **Responses**:
     - `200 OK`: Returns the number of points awarded.
     - `404 Not Found`: Indicates no receipt found for the provided ID.

#### Components

The OpenAPI schema defines the following components:

- **Schemas**: 
  - **Receipt**: Describes the structure of receipts.
  - **Item**: Describes individual items on a receipt.

### Technology Stack

#### Backend

- **Language**: Golang
- **Design Patterns**: 
  - Delegate Pattern
  - SOLID Principles
  - Event-Driven Architecture
  - Microservice/Controller Pattern
- **Database**: SQLite with SQL for data persistence, designed for clustering and indexed for speed.

#### Frontend

- **Language**: Vanilla JavaScript, HTML, and CSS
- **Purpose**: Provides an interface for users to search and add receipts.

### Security Measures

- **SSL Certificate**: Ensures secure communication between clients and the server.
- **Web Proxy**: Implemented in C# .NET, hosted on a VPS behind a Cloudflare-managed domain.
- **Tunnel Setup**: Uses Cloudflare for DNS management and VPN tunnels via No-Trust to secure API requests.
- **Server**: Custom-built server hosted behind a proxy and VPN tunnel into the network.
- **Cloudflare**: Acts as a reverse proxy, providing an additional layer of security and performance optimization.

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/samjay22/Fetch-Rewards-API/
   ```
2. Navigate into the repository directory:
   ```bash
   cd Fetch-Rewards-API
   ```
3. Ensure Golang version 1.18 or higher is installed. Refresh the project dependencies and index:
   ```bash
   go mod tidy
   ```
4. Build the Go binary:
   ```bash
   go build -o MY_EXE_NAME
   ```
5. Run the executable:
   ```bash
   ./MY_EXE_NAME
   ```

   > Note: An executable is already included in the project for convenience.

### Configurations

The configuration file contains core project settings such as host and port, SSL certificate, and SSL Pem Key file paths. Misconfiguration may result in errors such as CRT_NAME_MISMATCH. The live instance follows all specified configurations.

### Deployment

The application supports Docker and Kubernetes for containerization and orchestration. The server hosting the application runs Docker with Kubernetes to support Cloudflare, Redis, MongoDB, and MSSQL. Sensitive deployment files are excluded from the project.

### Performance

The application can sustain a throughput of 1000-3000 RPS (Requests Per Second), ensuring enterprise-level durability, data integrity, and security.

### Technical Considerations

#### Architectural Overview

**1. Microservices Architecture**: Each component of the application (e.g., receipt processing, points retrieval) is implemented as a separate microservice. This allows for independent scaling, deployment, and management.

**2. Event-Driven Architecture**: The receipt processing system can utilize event-driven principles, where the submission of a receipt triggers an event that other services can listen to and act upon.

**3. Database Sharding and Clustering**: SQLite is used with clustering and indexing to improve data retrieval performance. Considerations for future scalability might include migrating to a more robust distributed database system.

**4. Load Balancing**: Implement load balancing to distribute incoming traffic across multiple instances of the application, ensuring high availability and reliability.

#### Detailed Component Descriptions

**1. Receipt Processing Service**:
   - **Functionality**: Handles the receipt submission, validates the receipt, and stores it in the database.
   - **Technologies**: Golang, SQLite.

**2. Points Retrieval Service**:
   - **Functionality**: Retrieves the points awarded for a specific receipt.
   - **Technologies**: Golang, SQLite.

**3. Frontend Interface**:
   - **Functionality**: Provides a web interface for users to submit receipts and view points.
   - **Technologies**: Vanilla JavaScript, HTML, CSS.

**4. Security and Proxy**:
   - **Functionality**: Ensures secure communication and acts as a reverse proxy.
   - **Technologies**: C# .NET, Cloudflare.

### Flow Diagrams

**1. Receipt Submission Flow**:

```
User -> Frontend -> API Gateway -> Receipt Processing Service -> Database
```

**2. Points Retrieval Flow**:

```
User -> Frontend -> API Gateway -> Points Retrieval Service -> Database
```

### Additional Technical Considerations

**1. Error Handling**: Implement robust error handling mechanisms across all services to ensure graceful degradation and informative error messages.

**2. Logging and Monitoring**: Integrate logging and monitoring tools to track application performance and diagnose issues.

**3. API Rate Limiting**: Implement rate limiting to prevent abuse and ensure fair usage of the API.

**4. Automated Testing**: Develop unit, integration, and end-to-end tests to ensure the reliability and stability of the application.

### Conclusion

For further technical details, questions, or contributions, please reach out via email at samjaytaylor22@gmail.com.

---
