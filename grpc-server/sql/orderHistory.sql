-- Create the OrderHistoryLog table
CREATE TABLE OrderHistoryLog (
    id UUID PRIMARY KEY,
    spread FLOAT,
    aexcid VARCHAR(255),
    aprice FLOAT,
    apricevet FLOAT,
    avolume FLOAT,
    bexcid VARCHAR(255),
    bprice FLOAT,
    bpricevet FLOAT,
    bvolume FLOAT,
    createdAt TIMESTAMP
);