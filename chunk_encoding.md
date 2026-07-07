# Chunk Encoding

Http mechanism to breack the data response into smaller independent pieces.

Requires the `Transer-Encoding: chunked` header from the server/client(and the content length is not used)
Structure: 
- Size Indicator: each chunk has a length in hex and a `\r\n`
- Chunk data: the payload of the size
- End: zero length chunk for showing the end of the data 

## Trailer
- headers at the end
- need to specify on the header about the trailers before the body/chunks(thus headers are async)


## Format
```bash
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
... repeat ...
0\r\n
\r\n
```

where <n> is the number of bytes in the chunk, in hex

- Streaming large amounts of data
- real-time updates
- sending data of unknown size

