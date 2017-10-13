import 'rc-calendar/assets/index.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Calendar from 'rc-calendar';
import logo from './logo.svg';
import './App.css';

import { Button, Grid, Icon, Label, Menu, Table } from 'semantic-ui-react'

const TableDef = () => (
  <Table celled striped structured>
    <Table.Header>
      <Table.Row>
        <Table.HeaderCell width={2} />
        <Table.HeaderCell>Study Space 1</Table.HeaderCell>
        <Table.HeaderCell>Study Space 2</Table.HeaderCell>
        <Table.HeaderCell>Study Space 3</Table.HeaderCell>
        <Table.HeaderCell>Study Space 4</Table.HeaderCell>
      </Table.Row>
    </Table.Header>

    <Table.Body>
      <Table.Row>
        <Table.Cell rowSpan='2'>12:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>1:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>2:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>3:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>4:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>5:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>6:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>7:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>8:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>9:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>10:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>11:00 AM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>12:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>1:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>2:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>3:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>4:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>5:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>6:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>7:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>8:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>9:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>10:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell rowSpan='2'>11:00 PM</Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
      <Table.Row>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
        <Table.Cell selectable></Table.Cell>
      </Table.Row>
    </Table.Body>
  </Table>
)

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Felis Resource Management System</h1>
        </header>
        <Button>Login</Button>
        <Menu floated='right' pagination>
        <Menu.Item as='a' icon>
          <Icon name='left chevron' />
        </Menu.Item>
        <Menu.Item as='a'>1</Menu.Item>
        <Menu.Item as='a'>2</Menu.Item>
        <Menu.Item as='a'>3</Menu.Item>
        <Menu.Item as='a'>4</Menu.Item>
        <Menu.Item as='a' icon>
          <Icon name='right chevron' />
        </Menu.Item>
      </Menu>
        {/*<p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
         </p>*/}

        <Grid celled>
          <Grid.Row>
            <Grid.Column width={3}>
              <p>Find available resource today.</p>
              <Calendar />
            </Grid.Column>
            <Grid.Column width={13}>
              <TableDef />
            </Grid.Column>
          </Grid.Row>

        </Grid>
      </div>
    );
  }
}

export default App;
