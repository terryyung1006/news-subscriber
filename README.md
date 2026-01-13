# News Subscriber

This project aims to provide a news subscription service so that users can describe what kind of news they are interested in, and we will feed a summarized report to them everyday.

## 🚀 Quick Start

- **System Architecture** → See [`ARCHITECTURE.md`](./ARCHITECTURE.md)
- **Backend Setup** → See [`news-subscriber-core/README.md`](./news-subscriber-core/README.md)
- **Frontend Setup** → See [`news-subscriber-frontend/README.md`](./news-subscriber-frontend/README.md)
- **Feature Specs** → Each service has a `spec/` folder for feature documentation

## Architecture

This project has several parts:

- **news-fetcher**: Runs cronjobs to query the latest news from different sources, puts it into a standardized structure, and produces it into a Kafka topic. It is designed to be stateless.
- **etl-worker**: Consumes messages from Kafka, converts them into the desired format, and saves them to different databases like Chroma, PostgreSQL, and ClickHouse.
- **news-subscriber-frontend**: A web app for users to login, describe what kind of news they are interested in, and read their daily report.
- **news-subscriber-core**: The core engine of the system. It has access to all databases, processes all requests from the frontend, and handles different business logic.

## Development Guidelines

- **Package Selection**: If you need to import a package or library, choose the one that is most popular and actively maintained.
- **Design Patterns**: Follow design pattern rules so that the code is easy to maintain, open for extension, and closed for modification.
- **Simplicity**: Keep it simple and make sure it can run first. Do not create too many files.
- **Organization**: Keep the hierarchy simple and well-organized.
