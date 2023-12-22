// src/App.js
import React, { useState } from 'react';
import { Container, Row, Col, Navbar, Nav } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import NodeNetworkPage from './NodeNetworkPage';
import NodeTaskPage from './NodeTaskPage';
import PeerNodeDataPage from './PeerNodeDataPage';
import DataPage from './DataPage';

const App = () => {
  const [currentPage, setCurrentPage] = useState('nodeNetwork'); // Default page

  const renderPage = () => {
    switch (currentPage) {
      case 'nodeTask':
        return <NodeTaskPage />;
      case 'peerNodeData':
        return <PeerNodeDataPage />;
      case 'data':
        return <DataPage />;
      default:
        return <NodeNetworkPage/>
        
    }
  };

  return (
    <Container fluid>
      <Navbar bg="light" expand="lg">
        <Navbar.Brand>Node App</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="mr-auto">
            <Nav.Link onClick={() => setCurrentPage('nodeNetwork')}>Node Network</Nav.Link>
            <Nav.Link onClick={() => setCurrentPage('nodeTask')}>Node Task</Nav.Link>
            <Nav.Link onClick={() => setCurrentPage('peerNodeData')}>Peer Node Data</Nav.Link>
          
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    <div>
          {renderPage()}
          </div>
       
      </Container>
  );
};

export default App;
