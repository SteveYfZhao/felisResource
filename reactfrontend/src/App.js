import 'rc-calendar/assets/index.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Calendar from 'rc-calendar';
import logo from './logo.svg';
import './App.css';
import { Button, Grid, Icon, Label, Menu, Table } from 'semantic-ui-react';

var moment = require('moment');


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
                
                <TGrid/>
              </Grid.Column>
            </Grid.Row>

          </Grid>
        </div>
      );
    }
  }
  class TCell extends React.Component {
    render() {
      if (this.props.SpanOverride) {
        return;
      } else {
        return <Table.Cell rowSpan={this.props.rowspan} selectable={this.props.selectable}></Table.Cell>;
      }
    }
  }

  class THeader extends React.Component {
    render() {
      //this.props: rlist(array), colPerPage(int), pageNum(int)
      return <Table.Header>
        <Table.Row>
          <Table.HeaderCell width={2} />
          <Table.HeaderCell>Study Space 1</Table.HeaderCell>
          <Table.HeaderCell>Study Space 2</Table.HeaderCell>
          <Table.HeaderCell>Study Space 3</Table.HeaderCell>
          <Table.HeaderCell>Study Space 4</Table.HeaderCell>
        </Table.Row>
      </Table.Header>;
    }
  }

  class TRowSlot extends React.Component {
    render() {
      //this.props: rlist(array), colPerPage(int), pageNum(int)
      var cells = [];
      var rlist = this.props.rlist;
      var colPerPage = this.props.colPerPage;
      var pageNum = this.props.pageNum;
      if (!rlist[0].SpanOverride){
        cells.push(<Table.Cell rowSpan = {rlist[0].rowSpan} selectable = {rlist[0].selectable}>{rlist[0].value}</Table.Cell>);
      }
      for (var i=0; i < this.props.colPerPage; i++) {
        var cellData = rlist[i + 1 + pageNum * colPerPage];
          if (cellData == null || cellData.SpanOverride){
            continue;
          } else {
            console.log("push cell", cells);
            cells.push(<Table.Cell rowSpan = {cellData.rowSpan} selectable = {cellData.selectable}>{cellData.value}</Table.Cell>);
            
          }          
      }

     
      return <Table.Row>{cells}</Table.Row>;
    }
  }

  class TGrid extends React.Component {
    render() {
      //this.props: rlist(array), colPerPage(int), pageNum(int)
      const hourRow = 24;
      let colPerPage = 4;
      let pageNum = 0;
      let rMatrix = [];
      let rMatrixJSX = [];
      let momentParser = moment("2015-01-01").startOf('day'); //Set a fixed date to avoid leap year and daylightsaving issues.
      let totalItem = 12; // for test and debug only

      for (let h = 0 ; h < hourRow * 2; h++){
        let row = [];
        let hourCell = {rowSpan:1, selectable:false, SpanOverride:false, value:""};
        if ( h % 2 == 0 ){
          row.push({rowSpan:2, selectable:false, SpanOverride:false, value:momentParser.format("hh:mm a")});
        } else if (h % 2 == 1) {
          row.push({rowSpan:1, selectable:false, SpanOverride:true, value:""});
        }        
        momentParser.add(30, "m");

        for (let i = 0; i < totalItem; i++){
          var cellData = {rowSpan:1, selectable:true, SpanOverride:false, value:""};          
          row.push(cellData);
        }
        console.log("datarow", row);
        rMatrixJSX.push(<TRowSlot id={"hour"+h} rlist={row} colPerPage = {colPerPage} pageNum = {pageNum}/>)
        //rMatrixJSX.push(<Table.Row><Table.Cell rowSpan = '1' selectable = 'true'></Table.Cell></Table.Row>);
        
      }
      return <Table celled striped structured>
        <THeader/>
        <Table.Body>
        {rMatrixJSX}
        </Table.Body>
        
      </Table>;
    }
  }

  export default App;
