// src/DataPage.js
import React, { useState } from 'react';
import { Table, Button } from 'react-bootstrap';

const DataPage = () => {
  const [data, setData] = useState([
    { dataId: 1, url: 'https://example.com/data1', source: 'Node 1' },
    { dataId: 2, url: 'https://example.com/data2', source: 'Node 2' },
    { dataId: 3, url: 'https://example.com/data3', source: 'Node 3' },
  ]);

  const handleDownload = (dataId) => {
    // Add your download logic here
    console.log(`Downloading data with ID ${dataId}`);
  };

  return (
    <div>
      <h1>Data Page</h1>
      <Table striped bordered hover>
        <thead>
          <tr>
            <th>Data ID</th>
            <th>URL</th>
            <th>Source</th>
            <th>Download</th>
          </tr>
        </thead>
        <tbody>
          {data.map((row) => (
            <tr key={row.dataId}>
              <td>{row.dataId}</td>
              <td>{row.url}</td>
              <td>{row.source}</td>
              <td>
                <Button variant="success" onClick={() => handleDownload(row.dataId)}>
                  Download
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  );
};

export default DataPage;
