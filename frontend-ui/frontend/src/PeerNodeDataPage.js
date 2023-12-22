// src/PeerNodeDataPage.js
import React, { useState, useEffect } from 'react';
import { Form, Table, Button, Dropdown } from 'react-bootstrap';
import axios from 'axios';

const PeerNodeDataPage = () => {
  const [peerNodes, setPeerNodes] = useState([]);
  const [selectedPeerNode, setSelectedPeerNode] = useState(null);
  const [nodeData, setNodeData] = useState([]);

  useEffect(() => {
    const fetchPeerNodes = async () => {
      try {
        const response = await axios.get(`http://${window.url}/nodes`);
        setPeerNodes(response.data);
      } catch (error) {
        console.error('Error fetching peer nodes:', error);
      }
    };

    fetchPeerNodes();
  }, []);

  const handlePeerNodeChange = async (peerNodeId) => {
    try {
      const selectedNode = peerNodes.find((node) => node.peerid === peerNodeId);
      setSelectedPeerNode(selectedNode);

      const response = await axios.post(`http://${window.url}/nodes/peerdata`, {
        id: peerNodeId,
      });

      setNodeData(response.data);
    } catch (error) {
      console.error('Error fetching peer data:', error);
    }
  };

  const handleSyncData = async (dataIdToSync) => {
    try {
      const response = await axios.post(`http://${window.url}/collection/req`, {
        id: dataIdToSync,
        peerid: selectedPeerNode.peerid,
      });
  
      if (response.data.success) {
        console.log(`Sync successful for Data ID ${dataIdToSync}`);
      } else {
        console.error(`Sync failed for Data ID ${dataIdToSync}`);
      }
    } catch (error) {
      console.error('Error syncing data:', error);
    }
  };

  return (
    <div>
      <h1>Peer Node Data</h1>
      <Form>
        <Form.Group controlId="peerNode">
          <Form.Label>Select Peer Node</Form.Label>
          <Dropdown onSelect={(selectedNodeId) => handlePeerNodeChange(selectedNodeId)}>
            <Dropdown.Toggle variant="success" id="dropdown-basic">
              {selectedPeerNode ? selectedPeerNode.peerid : 'Select Peer Node'}
            </Dropdown.Toggle>

            <Dropdown.Menu>
              {peerNodes.map((node) => (
                <Dropdown.Item key={node.peerid} eventKey={node.peerid}>
                  {node.peerid}
                </Dropdown.Item>
              ))}
            </Dropdown.Menu>
          </Dropdown>
        </Form.Group>
      </Form>

      {selectedPeerNode && (
        <div>
          <h2>Data held by {selectedPeerNode.peerid}</h2>
          <Table striped bordered hover>
            <thead>
              <tr>
                <th>Data ID</th>
                <th>URL</th>
                <th>Description</th>
                <th>Interval</th>
                <th>Sync</th>
              </tr>
            </thead>
            <tbody>
              {nodeData.map((data) => (
                <tr key={data.DataID}>
                  <td>{data.DataID}</td>
                  <td>{data.Title}</td>
                  <td>{data.Description}</td>
                  <td>{data.Seconds}</td>
                  <td>
                    <Button variant="primary"onClick={() => handleSyncData(data.DataID)}>
                      Sync
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      )}
    </div>
  );
};

export default PeerNodeDataPage;
