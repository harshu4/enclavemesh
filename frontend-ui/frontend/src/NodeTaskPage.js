import React, { useState, useEffect } from 'react';
import { Table, Button, Form, Modal } from 'react-bootstrap';

const NodeTaskPage = () => {
  const [tasks, setTasks] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [newTaskDescription, setNewTaskDescription] = useState('');
  const [newTaskUrl, setNewTaskUrl] = useState('');
  const [newTaskInterval, setNewTaskInterval] = useState('');

  const handleStopTask = async (taskId) => {
    try {
      const removalData = { id: taskId }; // Construct the removal data

      const response = await fetch(`http://20.248.176.33:9900/nodes/rdata`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(removalData),
      });

      if (!response.ok) {
        throw new Error(`Failed to stop task: ${response.statusText}`);
      }

      // Update the state or perform any other necessary actions
      // For example, you can remove the stopped task from the tasks array
      setTasks((prevTasks) => prevTasks.filter((task) => task.DataID !== taskId));
    } catch (error) {
      console.error(error);
    }
  };

  const handleAddTask = async () => {
    try {
      console.log(JSON.stringify({
        description: newTaskDescription,
        url: newTaskUrl,
        interval: parseInt(newTaskInterval, 10),
      }))
      const response = await fetch('http://20.248.176.33:9900/nodes/data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          description: newTaskDescription,
          url: newTaskUrl,
          interval:newTaskInterval,
        }),
      });

      if (!response.ok) {
        throw new Error(`Failed to add task: ${response.statusText}`);
      }

      const newTask = await response.json();

      // Update the state with the new task
      console.log(newTask)
      setTasks((prevTasks) => [...prevTasks, newTask]);
      setShowModal(false); // Close the modal after adding the task
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    // Fetch the initial tasks from the API when the component mounts
    const fetchTasks = async () => {
      try {
        const response = await fetch('http://20.248.176.33:9900/nodes/getdata');
        
        if (response.ok) {
          let tasksData = await response.json();
          tasksData = tasksData.filter(item => item.Working);
          setTasks(tasksData);
        } else {
          throw new Error(`Failed to fetch tasks: ${response.statusText}`);
        }
      } catch (error) {
        console.error(error);
      }
    };

    fetchTasks();
  }, []); // Run once on component mount

  return (
    <div>
      <Button variant="primary" onClick={() => setShowModal(true)}>
        Add Task
      </Button>

      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Add Task</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form>
            <Form.Group controlId="taskDescription">
              <Form.Label>Description</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter task description"
                value={newTaskDescription}
                onChange={(e) => setNewTaskDescription(e.target.value)}
              />
            </Form.Group>
            <Form.Group controlId="taskUrl">
              <Form.Label>URL</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter task URL"
                value={newTaskUrl}
                onChange={(e) => setNewTaskUrl(e.target.value)}
              />
            </Form.Group>
            <Form.Group controlId="taskInterval">
              <Form.Label>Interval (in seconds)</Form.Label>
              <Form.Control
                type="number"
                placeholder="Enter task interval"
                value={newTaskInterval}
                onChange={(e) => setNewTaskInterval(e.target.value)}
              />
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowModal(false)}>
            Close
          </Button>
          <Button variant="primary" onClick={handleAddTask}>
            Add Task
          </Button>
        </Modal.Footer>
      </Modal>

      <Table striped bordered hover>
        <thead>
          <tr>
            <th>Description</th>
            <th>URL</th>
            <th>Interval (seconds)</th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((task) => (
            <tr key={task.Title}>
              <td>{task.Description}</td>
              <td>{task.Title}</td>
              <td>{task.Seconds}</td>
              <td>
                <Button variant="danger" onClick={() => handleStopTask(task.DataID)}>
                  Stop
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  );
};

export default NodeTaskPage;
