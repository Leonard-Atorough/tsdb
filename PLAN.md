We're building a simple time series database that we'll progressively build from a simple writer and reader to a more complex system with features like memory and disk indexing, compression, and query capabilities. The goal is to create a lightweight, efficient, and easy-to-use time series database that can handle high write throughput and provide fast read access to time series data.

We're going to approach this like a System design problem, starting with the basic requirements, building a simple naive implementation, and then iteratively improving it with more advanced features. We'll focus on the following key areas:
1. **Data Model**: Define how time series data will be structured and stored. This includes deciding on the schema, data types, and how to handle timestamps and values.
2. **Write Path**: Implement a simple writer that can accept time series data and store it in memory. We'll start with a basic in-memory structure and then explore options for persisting data to disk.
3. **Read Path**: Implement a simple reader that can retrieve time series data based on time ranges and other query parameters. We'll start with a basic in-memory retrieval mechanism and then explore more efficient indexing and query strategies.
4. **Indexing**: Explore different indexing strategies to improve read performance. This may include building in-memory indexes, disk-based indexes, and hybrid approaches that balance memory usage and query speed.
5. **Compression**: Implement data compression techniques to reduce storage requirements and improve read/write performance. We'll explore various compression algorithms and evaluate their trade-offs in terms of speed, compression ratio, and complexity.
6. **Query Capabilities**: Develop a query interface that allows users to retrieve time series data based on various parameters such as time ranges, aggregation functions, and filtering criteria. We'll start with basic query capabilities and then expand to support more complex queries and analytics.
7. **Scalability and Performance**: Evaluate the performance of the system under different workloads and optimize for high write throughput and low-latency reads. We'll explore techniques for scaling the system horizontally and vertically, as well as strategies for load balancing and fault tolerance.
8. **Testing and Benchmarking**: Implement a comprehensive testing framework to validate the correctness and performance of the system. We'll create benchmarks to measure write throughput, read latency, and query performance under various scenarios, and use these benchmarks to guide our optimization efforts.
9. **Documentation and User Experience**: Provide clear documentation and examples for users to understand how to use the time series database effectively. We'll focus on creating an intuitive API, comprehensive guides, and best practices for data ingestion, querying, and management.
By following this iterative approach, we aim to build a robust and efficient time series database that meets the needs of users while maintaining simplicity and ease of use. Each phase of development will be accompanied by thorough testing and performance evaluation to ensure that the system meets its design goals.

# Phase 1: Requirements and Basic Implementation
## Requirements
1. **Data Model**: The time series data will consist of a timestamp, a fieldset, a tagset and a measurement name. The timestamp will be represented as a Unix epoch time in milliseconds, the fieldset will contain the actual data points, and the tagset will contain metadata associated with the time series. The measurement name will be a string that identifies the type of data being stored.
2. **Write Path**: The writer will accept time series data in the form of (series_id, timestamp, value) tuples and store them in an in-memory data structure. The writer will support batch writes to improve throughput.
3. **Read Path**: The reader will allow users to query time series data based on series ID and time ranges. The reader will return the data in a structured format, such as a list of (timestamp, value) tuples.
4. **Indexing**: Initially, we will use a simple in-memory index to map series IDs to their corresponding time series data. This index will allow for efficient retrieval of data based on series ID. This will require understanding the best data structures for indexing, such as hash maps or trees, and implementing them in a way that supports fast lookups and insertions.

## Data Model
The time series data will be represented using the following struct:

```go
type TimeSeriesData struct {
    Measurement string `json:"m"`
    TagSet      map[string]string `json:"tags,omitempty"`
    FieldSet    map[string]any `json:"fields"`
    Timestamp   int64 `json:"ts"`
}
```

The `Measurement` field will store the name of the measurement, the `TagSet` will store metadata associated with the time series, the `FieldSet` will store the actual data points, and the `Timestamp` will store the timestamp of the data point.

The in-memory data structure for storing time series data will be a map that maps series IDs to slices of `TimeSeriesData` structs. This will allow for efficient retrieval of data based on series ID and time ranges.


## Files needed for Phase 1
- `internal/models/data.go`: This file will contain the definition of the `TimeSeriesData` struct and any other data models needed for the time series database.
- `internal/writer.go`: This file will contain the implementation of the writer that accepts time series data and stores it in memory or writes it to disk. The writer will support batch writes and will handle any necessary data validation and error handling.
- `internal/collector.go`: This file will contain the implementation of the collector that retrieves time series data based on series ID and time ranges. The collector will support filtering and aggregation of data, and will handle any necessary data validation and error handling.


## Phase 1 System design
The system will be design as a simple file writer that writes a data point to a file in JSON format. We'll expose a simple SDK that allows users to write data points 