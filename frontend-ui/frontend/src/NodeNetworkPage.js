// src/NodeNetworkPage.js
import React, { useState, useEffect } from 'react';
import ForceGraph3D from 'react-force-graph-3d';
import { Container, Row, Col, Form, Button, Spinner } from 'react-bootstrap';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';

const NodeNetworkPage = () => {
  const [data, setData] = useState({
    nodes: [
      { id: 'default-node', name: 'Default Node' }, // Default node
    ],
    links: [],
  });

  const [newNodeName, setNewNodeName] = useState('');
  const [newNodeUrl, setNewNodeUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const deleteRedundantNodes = (nodes, links, defaultNodeId) => {
    const connectedNodes = new Set();
    const connectedLinks = links.filter((link) => {
      connectedNodes.add(link.source);
      connectedNodes.add(link.target);
      return link.source === defaultNodeId || link.target === defaultNodeId;
    });
  
    const redundantNodes = nodes.filter((node) => !connectedNodes.has(node.id) && node.id !== defaultNodeId);
  
    return {
      nodes: nodes.filter((node) => !redundantNodes.includes(node)),
      links: connectedLinks,
    };
  };
  useEffect(() => {
    // Fetch nodes from the API
    const fetchNodes = async () => {
      try {
        const response = await axios.get(`http://${window.url}/nodes`); // Replace 'your_api_endpoint' with the actual API endpoint
        const apiNodes = response.data;

        // Create links from default node to each API node
        const apiLinks = apiNodes.map((apiNode) => ({
          source: 'default-node',
          target: apiNode.peerid.toString(),
          value: 10, // Adjust link value as needed
        }));

        // Update the data state with the fetched nodes and links
        console.log("am i called twice")
        const updatedData = deleteRedundantNodes([...data.nodes, ...apiNodes.map((apiNode) => ({ id: apiNode.peerid.toString(), name: apiNode.peerid.toString() }))],  [...data.links, ...apiLinks], 'default-node');
        setData((prevData) => (updatedData));
      } catch (error) {
        console.error('Error fetching nodes:', error);
      }
    };

    // Call the fetchNodes function
    fetchNodes();
  }, []); // Empty dependency array ensures the effect runs once when the component mounts

  const handleAddNode = async () => {
    try {
        setLoading(true)
      const response = await fetch(`http://${window.url}/nodes/add`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            id:newNodeUrl
        }),
      });

      if (!response.ok) {
        throw new Error(`Failed to add node: ${response.statusText}`);
      }

      const newNode = await response.json();
      
      setData((prevData) => ({
        nodes: [...prevData.nodes, { id: newNode.peerid.toString(), name: newNode.peerid.toString() }],
        links: [...prevData.links, { source: 'default-node', target: newNode.peerid.toString() }],
      }));
      // Update the state with the new node and a link connecting it to the default node
      setNewNodeName("")
      setNewNodeUrl("")
      setLoading(false)
    } catch (error) {
      console.error(error);
    }
  };


  return (
    <Container fluid style={{ position: 'relative' }}>
      <Row>
        <Col md={8} style={{ height: '100vh', position: 'relative', zIndex: 1 }}>
          <ForceGraph3D
            graphData={data}
            nodeAutoColorBy="name"
            linkColor={() => 'rgba(255, 255, 255, 10)'}
            onNodeClick={(node) => {
              console.log('Node clicked:', node);
            }}
            onNodeHover={(node) => {
              console.log('Node Hovered:', node);
            }}
            linkDirectionalParticles={(link) => link.value} // Set particle count based on link value
            linkDirectionalParticleSpeed={0.002} // Adjust particle speed
            linkWidth={2}
          />
        </Col>
        <Col md={3} style={{ alignContent: 'center', position: 'relative', zIndex: 2, right: '30%', textAlign: 'center', height: '40px' }}>
          <Form>
    
            <Form.Group controlId="nodeUrl">
              <Form.Label>Node URL</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter node URL"
                value={newNodeUrl}
                onChange={(e) => setNewNodeUrl(e.target.value)}
              />
            </Form.Group>
            <Button variant="primary" onClick={handleAddNode} disabled={loading} style={{ marginTop: '20px' }}>
              {loading ? (
                <>
                  <Spinner animation="border" size="sm" />
                  {' Adding Node...'}
                </>
              ) : (
                'Add Node'
              )}
            </Button>
          </Form>
        </Col>
      </Row>
    </Container>
  );
};

export default NodeNetworkPage;
