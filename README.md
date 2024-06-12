# jobet

Job scraper implementation, migrated from [goscrape](https://github.com/maxmwang/goscrape).

New features:
- Priority-rated companies: Higher priority companies are scrapped more frequently, while lower priority are scrapped less frequently. This is implemented to reduce outbound request rate.
- Pub-sub output: [ZeroMQ](https://zeromq.org/) is the primary form of output for scrape results, decoupling all handlers from the scraping daemon. 

## Technologies

- [SQLite](https://www.sqlite.org/): Lightweight SQL database
- [gRPC](https://grpc.io/): Lightweight communication between client and daemon 

