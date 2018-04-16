import 'rc-calendar/assets/index.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Calendar from 'rc-calendar';
import logo from './logo.svg';
import './App.css';
import axios from 'axios';
import { BrowserRouter, Link, Route } from 'react-router-dom'
import { Button, Grid, Icon, Label, Menu, Table, Input, Divider, Select, Checkbox, Sidebar, Segment, Header, Image, Form, List, Dropdown, Modal } from 'semantic-ui-react';

var moment = require('moment');
const hourInput = () => (
  <Input placeholder='1' />
)
const fromhourInput = () => (
  <Input placeholder='12' />
)
const apoptions = [
  { key: 'AM', text: 'AM', value: 'AM' },
  { key: 'PM', text: 'PM', value: 'PM' },
]
const halfhouroptions = [
  { key: '00', text: '00', value: '00' },
  { key: '30', text: '30', value: '30' },
]
const serverPortNum = 8081;

const serverProtocol = "http";
const pageSize = 100;

  class App extends Component {
    render() {
      return (
        <div className="App">
        {/*
          <header className="App-header">
            /*<img src={logo} className="App-logo" alt="logo" />
            <h1 className="App-title"><a href="/">Felis Resource Management System</a></h1>
          </header>
        */}
          <Route exact path="/" component={MainView}/>
          <Route path="/index" component={MainView}/> 
          <Route path="/adminview" component={AdminView}/>     
          <Route path="/login" component={LoginView}/>     
          <Route path="/signup" component={SignupView}/>     
        </div>
      );
    }
  }

  class MainView extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        resp: null,
        date:moment().format("YYYY-MM-DD"),
        startTime:null,
        startHr:moment().format("h"),
        startMin:"00",
        startAmPm:moment().format("A"),
        timeSpaninHr:1,
        viewMode:"day",
        bookingStatus:null,
      };
      this.handleInputChange = this.handleInputChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
      this.handleCalChange = this.handleCalChange.bind(this);
      this.getResBookingStatus = this.getResBookingStatus.bind(this);
      this.changeViewMode = this.changeViewMode.bind(this);
    }

    handleInputChange(event) {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;

      this.setState({
        [name]: value
      });
    }
  
    handleSubmit(event) {
      alert('A name was submitted: ' + this.state.value);
      event.preventDefault();
    }

    handleCalChange(value) {
      var datestr = moment(value._d).format("YYYY-MM-DD");
      console.log('on panel change', value._d, datestr);
      this.setState({
        ["date"]: datestr
      });
    }

    handleDropdownChange = (e, target) => {
      this.setState({
        [target.name]: target.value
      });
    }

    searchAvailRoom(){
      var startTimeTemp = this.state.date + " " + this.state.startHr + ":" + this.state.startMin ;
      var startMmt = moment(startTimeTemp);
      if (this.state.startAmPm == "PM") {
        startMmt.add(12,'h')
      }
      var startTime = startMmt.format();
      
      var endTime = startMmt.add(this.state.timeSpaninHr, 'h').format();
      console.log("startMmt", startMmt, startMmt.add(this.state.timeSpaninHr, 'h'), startMmt.add(this.state.timeSpaninHr, 'h').format("YYYY-MM-DD, h:mm a"))



      var resp = null;
      var self = this;

      axios.post(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/getbookableresforuserattime', {
          startTime   :startTime,
          endTime       :endTime,
        }, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          resp: response.data.Data
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 
    }

    getResBookingStatus(){
      var startTimeTemp = this.state.date;
      var startMmt = moment(startTimeTemp).startOf(this.state.viewMode);
      var endMmt = moment(startTimeTemp).endOf(this.state.viewMode);
      
      var startTime = startMmt.format();      
      var endTime = endMmt.format();
      
      console.log("start/end Time", startTime, endTime)

      var resp = null;
      var self = this;
      var rooms = "";
      this.state.resp.forEach(element => {
        if (rooms.length>0) {
          rooms = rooms + "," + element.Id
        } else {
          rooms = "" + element.Id
        }        
      });

      axios.post(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/getresbookingstatus', {
          rooms       :rooms,
          startTime   :startTime,
          endTime     :endTime,
        }, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          bookingStatus: response.data.Data
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 
    }

    changeViewMode(mode){
      this.setState({
        viewMode:mode
      });
      this.getResBookingStatus();
    }


    
    render() { 
      let listItems = [];
      let listArea = null;
      if (this.state.resp != null && Array.isArray(this.state.resp)) {
        listArea = <List divided relaxed className="item-list user-list"> {listItems} </List>
        var self = this;
        if (this.state.resp && Array.isArray(this.state.resp)){
          this.state.resp.forEach(function(element){
            listItems.push(
              <List.Item key={element.Id}>
                <List.Icon name='user outline' size='large' verticalAlign='middle' />
                <List.Content>
                  <List.Header>
                  </List.Header>
                  <List.Description>                                  
                      <span>Room: {element.Name}</span><br/>  
                      <Button className="editUserBtn" >View room schedule</Button>
                  </List.Description>
                </List.Content>
              </List.Item>
            );
          });
        }        
      } else {
        listArea = <p>Please search or browse users.</p>
      }
      return (
        <div className="MainView">
                  {/*<p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
         </p>*/}
          
          {/*
          <LoginCtrls></LoginCtrls>
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


          <Grid celled>
            <Grid.Row>
              <Grid.Column width={3}>
                <p>Select date</p>
                <Calendar mode="date" format="YYYY-MM-DD" onChange={this.handleCalChange}/>
                
                <div>
                  <h3>Search available rooms</h3>
                  <p>选中日期和时间后，右侧显示的时间也发生变化，对应的时间槽出现高亮条。（待定）Y坐标可以在小时和一周七天和一月之间切换。Day/Week/Month view。需要从服务器返回2D Array[room][booking info]</p>
                  <div>
                    <div>From:  
                      <Input type='text' size='medium'  name="startHr" value={this.state.startHr}  onChange={this.handleInputChange} placeholder='12' action>
                      <input size = '2'/>                  
                      </Input> : 
                      <Select compact name="startMin" value={this.state.startMin}  onChange={this.handleDropdownChange} options={halfhouroptions}  />
                      <Select compact name="startAmPm" value={this.state.startAmPm} onChange={this.handleDropdownChange} options={apoptions}  /></div>
                      <p>for at least <input name="timeSpaninHr" value={this.state.timeSpaninHr} onChange={this.handleInputChange} placeholder='1' size = '2'/> hour</p>
                      <Button onClick={() =>this.searchAvailRoom()}>Search</Button>
                    </div>

                  <Divider horizontal>Or</Divider>
                  <Button>Show all</Button>
                  <Divider horizontal>Switch View</Divider>
                  <Button onClick={()=>this.changeViewMode("day")}>Day</Button>
                  <Button onClick={()=>this.changeViewMode("week")}>Week</Button>
                  <Button onClick={()=>this.changeViewMode("month")}>Month</Button>
                  {listArea}
                </div>
              </Grid.Column>
              <Grid.Column width={13}>
                
                <TGrid bookingdata={this.state.bookingStatus}/>
              </Grid.Column>
            </Grid.Row>

          </Grid>
          */}

          <div className="header">
            <LoginCtrls/>
            <div>
                    <div>From:  
                      <Input type='text' size='medium'  name="startHr" value={this.state.startHr}  onChange={this.handleInputChange} placeholder='12' action>
                      <input size = '2'/>                  
                      </Input> : 
                      <Select compact name="startMin" value={this.state.startMin}  onChange={this.handleDropdownChange} options={halfhouroptions}  />
                      <Select compact name="startAmPm" value={this.state.startAmPm} onChange={this.handleDropdownChange} options={apoptions}  />
                      <p>for at least <input name="timeSpaninHr" value={this.state.timeSpaninHr} onChange={this.handleInputChange} placeholder='1' size = '2'/> hour</p>
                      <Button onClick={() =>this.searchAvailRoom()}>Search</Button>  <Button>Show all</Button>
                    </div>

                  <Button onClick={()=>this.changeViewMode("day")}>Day</Button>
                  <Button onClick={()=>this.changeViewMode("week")}>Week</Button>
                  <Button onClick={()=>this.changeViewMode("month")}>Month</Button>
                </div>

          </div>
         
          <FlGrid bookingdata={this.state.bookingStatus}/>






        </div>
      );
    }

  }

  class FlGrid extends React.Component {
    constructor (props) {
      super(props);
      this.state = {
        bookingdata: '',
        modalOpen: false,
        selectedStartTime:"",
        selectedroom:"",
        selectedEndTIme:"",
      };
    }
    componentDidMount () {
      this.setState({bookingdata: this.props.bookingdata});
    }

    componentWillReceiveProps(newProps) {
      this.setState({bookingdata: newProps.bookingdata});
    }

    handleOpen = (t,r) => {
      console.log("t,r",t,r);
      this.setState({ 
        selectedStartTime: t,
        selectedroom:r,
        modalOpen: true 
      })
    }

    handleClose = () => this.setState({ modalOpen: false })


    render() {
      //this.props: rlist(array), colPerPage(int), pageNum(int)
      console.log("this.props.bookingdata",this.props.bookingdata); // after update, this logs the updated name
      console.log("this.state.bookingdata",this.state.bookingdata);
      const hourRow = 24;
      let colPerPage = 4;
      let pageNum = 0;
      let gridHeaders = [];
      let gridcolumns = [];
      let gridAxles = [];
      let gridHoriLns = [];
      let alldayevents = [];
      
      let momentParsera = moment("2015-01-01").startOf('day'); //Set a fixed date to avoid leap year and daylightsaving issues.
      let totalItem = 12; // for test and debug only

      gridHeaders.push(
        <div className="gridHeaderfiller" key="filler">&nbsp;</div>
      );

      if (this.state.bookingdata){
        var dict = this.state.bookingdata;
        for(var key in dict){
          gridHeaders.push(
            <div className="gridHeaderblock" key={dict[key].Id}>{dict[key].Name}</div>
          );
        }
      } else {
        for (let i = 0; i<colPerPage; i++) {
          gridHeaders.push(
            <div className="gridHeaderblock" key={i}>No room selected</div>
          );
        }
      }

      for (let a = 0 ; a < hourRow; a++){
        gridAxles.push(
          <div className="gridaxleblock" key={a}>{momentParsera.format("hh:mm a")}</div>
        );
        gridHoriLns.push(<div className="gridHoriLn" key={a}></div>);
        momentParsera.add(1, "h");
      }

      

      if (this.state.bookingdata){
        let dict = this.state.bookingdata;
        gridcolumns = [];       
        for(let key in dict){
          /*
          var cellData = {rowSpan:1, selectable:true, cellTime:momentParser.format("hh:mm a"), roomId:dict[key].Id,  SpanOverride:false, value:"h="+momentParser.format("hh:mm a")+" roomid= "+dict[key].Id};          
          row.push(cellData);
          */
          let gcCell = [];
          let momentParser30m = moment("2015-01-01").startOf('day'); //Set a fixed date to avoid leap year and daylightsaving issues.
          for (let h = 0 ; h < hourRow * 2; h++){
            gcCell.push(<div className="gridCell" key={dict[key].Id+""+h}>{dict[key].Id+" "+momentParser30m.format("hh:mm a")}</div>)
            momentParser30m.add(30, "m");
          }
          gridcolumns.push(<div className="gridCol" key={dict[key].Id}>{gcCell}</div>);
        }
      } else {
        /*
        for (let i = 0; i < totalItem; i++){
          
          var cellData = {
            rowSpan:1, 
            selectable:true, 
            SpanOverride:false, 
            //value:"h="+momentParser.format("hh:mm a")+" c= "+i
          };          
          row.push(cellData);            
        }
        */
       

        gridcolumns = [];
        for (let g = 0; g < colPerPage; g++) {
          let gcCell = [];
          let momentParser30m = moment("2015-01-01").startOf('day'); //Set a fixed date to avoid leap year and daylightsaving issues.
          for (let h = 0 ; h < hourRow * 2; h++){
            let t = momentParser30m.format("hh:mm a");
            gcCell.push(<div className="gridCell" key={g+""+h} onClick={() =>this.handleOpen(t, g)}>{g+" "+t}</div>)
            momentParser30m.add(30, "m");
          }
          gridcolumns.push(<div className="gridCol" key={g}>{gcCell}</div>);
        }
      }
      /*

      for (let h = 0 ; h < hourRow * 2; h++){
        
        let row = [];
        let hourCell = {rowSpan:1, selectable:false, SpanOverride:false, value:""};
        if ( h % 2 == 0 ){
          row.push({key:h, rowSpan:2, selectable:false, SpanOverride:false, value:momentParser.format("hh:mm a")});
        } else if (h % 2 == 1) {
          row.push({key:h, rowSpan:1, selectable:false, SpanOverride:true, value:"a"});
        }
        //momentParser.add(30, "m");

        //console.log("datarow", row);
        //rMatrixJSX.push(<TRowSlot key={"hour"+h} showbookingmodel={this.handleOpen} id={"hour"+h} rlist={row} colPerPage = {colPerPage} pageNum = {pageNum}/>)
        //rMatrixJSX.push(<Table.Row><Table.Cell rowSpan = '1' selectable = 'true'></Table.Cell></Table.Row>);
        
      }*/
      return (<div className="flbox">
        <div className="gridheader">
          <div className="gridheaderL1">
          {gridHeaders}            
          </div>          
        </div>
      <div className="gridbody">
        <div className="gridbodyL1">
          <div aria-hidden="true" className="gridHoriLns">
            {gridHoriLns}
          </div>
          
            <div className="gridaxle">
              <div className="gridaxleL1">
                {gridAxles}
              </div>            
            </div>
            <div className="gridscolumn">
              {gridcolumns}
            </div>
       
        </div>        
        
      </div>
        
        
        
        {/*
        <Table celled striped structured>
        <Table.Header>
        <Table.Row>
        <Table.HeaderCell width={2} />
        {rHeaders}
        </Table.Row>
        </Table.Header>
        <Table.Body>
        {rMatrixJSX}
        </Table.Body>        
      </Table>
        */}
      <Modal
        trigger={<Button onClick={this.handleOpen}>Show Modal</Button>}
        open={this.state.modalOpen}
        onClose={this.handleClose}        
        size='small'
      >
        <Header icon='browser' content='Cookies policy' />
        <Modal.Content>
          <h3>This website uses cookies to ensure the best user experience.</h3>
          <p>Start time:{this.state.selectedStartTime} </p>
          <p>Room:{this.state.selectedroom} </p>
          <input type='text' size='medium'  name="endTime" value={this.state.selectedEndTime}  onChange={this.handleInputChange} placeholder='12'/>
        </Modal.Content>
        <Modal.Actions>
          <Button color='green' onClick={this.handleClose} inverted>
            <Icon name='checkmark' /> Got it
          </Button>
          <Button color='grey' onClick={this.handleClose}>
            <Icon name='checkmark' /> Cancel
          </Button>
        </Modal.Actions>
      </Modal>

      </div>);
    }
  }


  class TRowSlot extends React.Component {
    constructor (props) {
      super(props);
      this.state = {
        
      };
      //this.handleOpen = this.handleOpen.bind(this);
    }

    handleOpen = (t,r) => {
      console.log("fire handleopen",t,r);

      this.props.showbookingmodel(t,r);
    };
    render() {
      //this.props: rlist(array), colPerPage(int), pageNum(int)
      let cells = [];
      let rlist = this.props.rlist;
      let rNum = this.props.rNum;
      let colPerPage = this.props.colPerPage;
      let pageNum = this.props.pageNum;

      //console.log("rlist", rlist);
      
      if (!rlist[0].SpanOverride){

        var key = rNum + "override";
        cells.push(<Table.Cell key={key} rowSpan = {rlist[0].rowSpan} selectable = {rlist[0].selectable}>{rlist[0].value}</Table.Cell>);
      }
      for (let i=0; i < this.props.colPerPage; i++) {
        let cellData = rlist[i + 1 + pageNum * colPerPage];
          if (cellData == null || cellData.SpanOverride){
            continue;
          } else {
            //console.log("push cell", cells);
            let key = rNum + "c" + i
            cells.push(<Table.Cell key = {key} onClick={()=>this.handleOpen(cellData.cellTime, cellData.roomId)} rowSpan = {cellData.rowSpan} selectable = {cellData.selectable}>{cellData.value}</Table.Cell>);
            
          }          
      }

     
      return <Table.Row>
      {cells}
      </Table.Row>;
    }
  }
/*
  class TGrid extends React.Component {
    constructor (props) {
      super(props);
      this.state = {
        bookingdata: '',
        modalOpen: false,
        selectedStartTime:"",
        selectedroom:"",
      };
    }
    componentDidMount () {
      this.setState({bookingdata: this.props.bookingdata});
    }

    componentWillReceiveProps(newProps) {
      this.setState({bookingdata: newProps.bookingdata});
    }

    handleOpen = (t,r) => {
      console.log("t,r",t,r);
      this.setState({ 
        selectedStartTime: t,
        selectedroom:r,
        modalOpen: true 
      })
    }

    handleClose = () => this.setState({ modalOpen: false })


    render() {
      //this.props: rlist(array), colPerPage(int), pageNum(int)
      console.log("this.props.bookingdata",this.props.bookingdata); // after update, this logs the updated name
      console.log("this.state.bookingdata",this.state.bookingdata);
      const hourRow = 24;
      let colPerPage = 4;
      let pageNum = 0;
      let rHeaders = [];
      let rMatrix = [];
      let rMatrixJSX = [];
      let momentParser = moment("2015-01-01").startOf('day'); //Set a fixed date to avoid leap year and daylightsaving issues.
      let totalItem = 12; // for test and debug only

      if (this.state.bookingdata){
        var dict = this.state.bookingdata;
        for(var key in dict){
          rHeaders.push(
            <Table.HeaderCell key={dict[key].Id}>{dict[key].Name}</Table.HeaderCell>
          );
        }
      } else {
        for (let i = 0; i<colPerPage; i++) {
          rHeaders.push(
            <Table.HeaderCell key={i}>&nbsp;</Table.HeaderCell>
          );
        }
      }

      for (let h = 0 ; h < hourRow * 2; h++){
        let row = [];
        let hourCell = {rowSpan:1, selectable:false, SpanOverride:false, value:""};
        if ( h % 2 == 0 ){
          row.push({key:h, rowSpan:2, selectable:false, SpanOverride:false, value:momentParser.format("hh:mm a")});
        } else if (h % 2 == 1) {
          row.push({key:h, rowSpan:1, selectable:false, SpanOverride:true, value:"a"});
        } 
        if (this.state.bookingdata){
          let dict = this.state.bookingdata;          
          for(let key in dict){
            let cellData = {rowSpan:1, selectable:true, cellTime:momentParser.format("hh:mm a"), roomId:dict[key].Id,  SpanOverride:false, value:"h="+momentParser.format("hh:mm a")+" roomid= "+dict[key].Id};          
            row.push(cellData);
          }
        } else {
          for (let i = 0; i < totalItem; i++){
            let cellData = {
              rowSpan:1, 
              selectable:true, 
              SpanOverride:false, 
              //value:"h="+momentParser.format("hh:mm a")+" c= "+i
            };          
            row.push(cellData);
          }

        }
        momentParser.add(30, "m");

        //console.log("datarow", row);
        rMatrixJSX.push(<TRowSlot key={"hour"+h} showbookingmodel={this.handleOpen} id={"hour"+h} rlist={row} colPerPage = {colPerPage} pageNum = {pageNum}/>)
        //rMatrixJSX.push(<Table.Row><Table.Cell rowSpan = '1' selectable = 'true'></Table.Cell></Table.Row>);
        
      }
      return (<div><Table celled striped structured>
        <Table.Header>
        <Table.Row>
        <Table.HeaderCell width={2} />
        {rHeaders}
        </Table.Row>
        </Table.Header>
        <Table.Body>
        {rMatrixJSX}
        </Table.Body>        
      </Table>
      <Modal
        trigger={<Button onClick={this.handleOpen}>Show Modal</Button>}
        open={this.state.modalOpen}
        onClose={this.handleClose}
        
        size='small'
      >
        <Header icon='browser' content='Cookies policy' />
        <Modal.Content>
          <h3>This website uses cookies to ensure the best user experience.</h3>
          <p>Start time:{this.state.selectedStartTime} </p>
          <p>Room:{this.state.selectedroom} </p>
        </Modal.Content>
        <Modal.Actions>
          <Button color='green' onClick={this.handleClose} inverted>
            <Icon name='checkmark' /> Got it
          </Button>
        </Modal.Actions>
      </Modal>

      </div>);
    }
  }
  */

  class LoginCtrls extends React.Component {
    constructor(props) {
      super(props);
      this.state = {resp: null};
    }

    logout = () => {
      document.cookie = "data=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      console.log("this", this);
      this.setState({
        resp: null
      });
    }

    getData = () => {
      var resp = null;
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/userbasicinfo', {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          resp: response.data.Data
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      });
    }

    componentDidMount() {
      this.getData();      
    }

      
    componentWillUnmount() {
      
    }

    render() {
     
      //console.log("resp.Username", resp.Username);
      if (this.state.resp && this.state.resp.Username) {
        return (
          <div>
            <p>Hello, {this.state.resp.Username}</p>
            <Button onClick={this.logout}>Log out</Button>
            <CtrlPanel isAdmin = {this.state.resp.CommonPermissions.indexOf("basicAdmin") > -1} />
          </div>
        );

      } else {
        return (
          <p>
            <Button to="/login" as = {Link}>login</Button>
            <Button  to = "/signup" as = {Link}>signup</Button>
          </p>
        );
      }

    }
  }

  class CtrlPanel extends React.Component {
    render() {
      if (this.props.isAdmin) {
        return (
          <div>
          <nav>
            {/*<Link to="/dashboard">Admin Control Panel</Link>*/}
            <Link to="/adminview">Admin Control Panel</Link>
          </nav>          
        </div>
          
        );
      } else {
        return <a>User Settings</a>
      }
    }
  }

  class AdminView extends React.Component {
    constructor(props) {
      super(props);  
      this.state = { 
        isAdmin: false,
      };
    }

    getData = () => {
      var resp = null;
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/userbasicinfo', {withCredentials: true})
      .then(function (response) {
        //console.log(response);
        if (response.data.Data.CommonPermissions.indexOf("basicAdmin") > -1) {
          self.setState({
            isAdmin: true
          });
        }        
        //console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      });
    }

    componentDidMount() {
      this.getData();      
    }

    render() {
      if (this.state.isAdmin) 
      {
        return (
          <div className="AdminView">
            <SidebarLeftSlideAlong/>



            {/*<Grid celled>
              <Grid.Row>
                <Grid.Column width={3}>
                  <div className = "sidebar">
                    <p>Sidebar</p>
            
                  </div>
                </Grid.Column>
                <Grid.Column width={13}>
                  <div className = "mainbody">
                    <h2>AdminView</h2>
                  </div>              
                </Grid.Column>
              </Grid.Row>

            </Grid>*/}
          </div>
          )

      } else {
        return (
          <div className="AdminView">
            You do not have permission to view this page.
          </div>
          )

      }
      
        
      
    }
  }

  class SidebarLeftSlideAlong extends Component {
    constructor(props) {
      super(props);  
      this.state = { 
        visible: true,
        activeItem: null
      };
    }

    handleItemClick = name => this.setState({ activeItem: name })
    toggleVisibility = () => this.setState({ visible: !this.state.visible })
    
  
    render() {
      const { activeItem } = this.state || {}
      const { visible } = this.state
      return (
        <div>
          
          <Sidebar.Pushable as={Segment}>
            <Sidebar as={Menu} animation='slide along' width='thin' visible={visible} icon='labeled' vertical inverted>

            {/*
            <Menu.Header name='rooms'>
              <Icon name='browser' />
              Manage Rooms
            </Menu.Header>
            <Menu.Menu>
              <Menu.Item name='Add room' active={activeItem === 'enterprise'} onClick={this.handleItemClick} />
              <Menu.Item name='consumer' active={activeItem === 'consumer'} onClick={this.handleItemClick} />
            </Menu.Menu>
            */}  


              <Menu.Item name='rooms' active={activeItem === 'rooms'} onClick={this.handleItemClick} as = {Link} to="/adminview/rooms">
               <Icon name='browser' /> Manage Rooms
              </Menu.Item>

              <Menu.Item name='resources' active={activeItem === 'resources'} onClick={this.handleItemClick} as = {Link} to="/adminview/resources">
                <Icon name='tags' /> Manage Resources and Tags
              </Menu.Item>
              <Menu.Item name='booking' active={activeItem === 'booking'} onClick={this.handleItemClick} as = {Link} to="/adminview/booking">
                <Icon name='tags' /> Manage booking
              </Menu.Item>
              <Menu.Item name='users' active={activeItem === 'users'} onClick={this.handleItemClick} as = {Link} to="/adminview/users">
                
                <Icon name='users' />
                Manage Users

              </Menu.Item>
              <Menu.Item name='stats' active={activeItem === 'stats'} onClick={this.handleItemClick} as = {Link} to="/adminview/stats">
                
                <Icon name='bar chart' />
                Stats

              </Menu.Item>
            </Sidebar>
            <Sidebar.Pusher>
              <Segment basic>
              <Header as='h3'><Button onClick={this.toggleVisibility}>Menu</Button> Application Content</Header>
              <Route exact path={'/adminview'} component={ManageRoom}/>
              <Route path={'/adminview/rooms'} component={ManageRoom}/>
              <Route path={'/adminview/resources'} component={ManageRes}/>
              <Route path={'/adminview/booking'} component={ManageBooking}/>
              <Route path={'/adminview/users'} component={ManageUser}/>
              <Route path={'/adminview/stats'} component={ShowStats}/>
              
                
                
                <br/><br/><br/><br/><br/><br/><br/>
                <Image src='/assets/images/wireframe/paragraph.png' />
                <br/><br/><br/><br/><br/><br/><br/>
                
              </Segment>
            </Sidebar.Pusher>
          </Sidebar.Pushable>
        </div>
      )
    }
  }

  class ManageRoomMain extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        resourceid: '',
        displayname: '',  
        restype: "room",
        viewpermission:'userperm',
        bookpermission:'userperm',
        capacity: 1,
        resp: null,
        resourcelist: null,
        internalonly:false,
      };
  
      this.handleInputChange = this.handleInputChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange(event) {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;

      //console.log("handle change: type:", event.type, "checked: ", target.checked, "value:", target.value );

      this.setState({
        [name]: value
      });
      if (name == "internalonly" && value) {
        this.setState({
          ['viewpermission']: "advuserperm",
          ['bookpermission']: "advuserperm",
        });
        console.log("Change to internal");
      } 
      if(name == "internalonly" && !value) {
        this.setState({
          ['viewpermission']: "userperm",
          ['bookpermission']: "userperm",
        });
        console.log("Change to normal user");
        
      }
    }
  
    handleSubmit(event) {
      alert('A name was submitted: ' + this.state.value);
      event.preventDefault();
    }
    addResource() {
      var resp = null;
      var self = this;

      axios.post(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/addresource', {

          displayname   :self.state.displayname,
          restype       :self.state.restype,
          viewpermission:self.state.viewpermission,
          bookpermission:self.state.bookpermission,
          capacity      :self.state.capacity,
          pageSize      :20
        }, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        /*
        self.setState({
          resp: response.data.Data
        });
        
        console.log("this.state.resp", self.state.resp);
        resp = response;
        */

      })
      .catch(function (error) {
        console.log(error);
      }); 

    }
    listResource(offset) {
      var resp = null;
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/listresourceforadm', {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          resp: response.data.Data
        });
        
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 
    }
    render() {
      

      let listItems = [];
      let listArea = null;
      if (this.state.resp != null) {
        listArea = <List divided relaxed className="item-list resource-list"> {listItems} </List>
        this.state.resp.forEach(function(element){
          listItems.push(
            <List.Item key={element.id}>
              <List.Icon name='user outline' size='large' verticalAlign='middle' />
              <List.Content>
                <List.Header>
                </List.Header>
                <List.Description>
                    <span>Name: {element.displayname}</span><br/>

                    <Button className="editUserBtn" as = {Link} to={"/adminview/rooms/detail/"+element.id}>Edit this room</Button>
                  
                </List.Description>
              </List.Content>
            </List.Item>
          );
        });
      } else {
        listArea = <p>Please search or browse rooms.</p>
      }
      
        return (
          <div className="ManageRooms">
            <p>Manage rooms</p>

            <p>Add room</p>
            <Form>
            {/*
            <Form.Field>
              <label>Unique ID</label>
              <input className="form-text-input" name="resourceid" value={this.state.userid} placeholder='User ID' onChange={this.handleInputChange}/>
            </Form.Field>
            */}
            <Form.Field>
              <label>Room name</label>
              <input className="form-text-input" name="displayname" value={this.state.displayname} placeholder='Room name' onChange={this.handleInputChange}/>
            </Form.Field>
            <Form.Field>
              <label>Capacity</label>
              <input className="form-text-input" name="capacity" value={this.state.capacity} placeholder='1' onChange={this.handleInputChange} />
            </Form.Field>
            <Form.Field>
              <label>Available to interal users only: </label>
                <input name="internalonly" type="checkbox" checked={this.state.internalonly} onChange={this.handleInputChange}/>
              
            </Form.Field>
            <Button type='submit' onClick={() =>this.addResource()}>Submit</Button>
          </Form>

          <br/>
          <p>Browse Rooms</p>
          <Button onClick={() =>this.listResource()}>Browse Rooms</Button>



            <p>Remove room</p>
            <p>Edit room</p>
            {listArea}

          </div>
          )      
    }
  }

  class RoomDetails extends React.Component {
    constructor(props) {
      super(props);  
      this.state = { 
        resp: null,
        lastRID: ""
      };
      this.handleInputChange = this.handleInputChange.bind(this);
      //this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange = (event) => {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;
      let updatedresp = Object.assign({}, this.state.resp);    //creating copy of object
      updatedresp[name] = value;                        //updating value

      console.log("onchange fired", event.target);
      console.log("onchange fired", event.target.checked);
      console.log("onchange fired", event.target.value);

      this.setState({
        resp: updatedresp
      });
    }
/*
    handleToggleChange = (event) => {
      const target = event.target;
      //const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.textContent;
      let updatedresp = Object.assign({}, this.state.resp);    //creating copy of object
      updatedresp[name] = !updatedresp[name];                        //updating value

      console.log("onchange fired", event.target);
      console.log("onchange fired", event.target.checked);


      this.setState({
        resp: updatedresp
      });
    }
    */
/*
    handleRoleInputChange = (event) => {
      const target = event.target;
      //const value = target.checked;
      const name = target.name;

      console.log("handleRoleInputChange fired", event.target);
      console.log("handleRoleInputChange fired", event.target.checked);
      console.log("handleRoleInputChange fired", event.target.value);

      let epName = (event.target.checked)? "assignroletouser" : "removerolefromuser";
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/'+epName+'?username='+this.props.match.params.uid+'&rolename=' + name, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        if (response.data.Data === "OK"){
          let updatedresp = Object.assign({}, self.state.resp);    //creating copy of object
          updatedresp.Roles[name] = !updatedresp.Roles[name];                        //updating value
          self.setState({
            resp: updatedresp
          });
        }
      })
      .catch(function (error) {
        console.log(error);
      });       
    }
  */

    updateResource = () => {
    var resp = null;
    var self = this;
    axios.post(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/editresource', {
        displayname   :self.state.resp.displayname,
        //restype       :self.state.resp.restype,
        restype       :"room",
        viewpermission:self.state.resp.viewpermission,
        bookpermission:self.state.resp.bookpermission,
        capacity      :self.state.resp.capacity,
        id            :self.state.resp.id,
      }, {withCredentials: true})
    .then(function (response) {
      console.log(response);
    })
    .catch(function (error) {
      console.log(error);
    }); 

    }
    handleCancel = () => {      
      this.getData();
    }

    getData = () => {
      var resp = null;
      let rID = this.props.match.params.rid
      var self = this;
      if (rID) {
        axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/fetchresourcedetailadm?rid='+rID, {withCredentials: true})
      .then(function (response) {
        console.log("response", response);
        if (response.data.Data && response.data.Data[0]){
          self.setState({
            resp: response.data.Data[0],
            lastRID:rID
          });
        }
        
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      });
      }      
    }
    componentDidMount() {
      this.getData()
    }

    componentDidUpdate(){
      if (this.props.match.params.rid != this.state.lastRID) {
        this.getData()
      }

    }

    render() {
      let rID = this.props.match.params.rid;
      let serverresp = this.state.resp;
      let output = <div></div>;
      let roleuis = [];

      if (serverresp!=null){

        output= <div className="RoomDetails">
        <Form>
          <Form.Field>
              <label>Room name</label>
              <input className="form-text-input" name="displayname" value={this.state.resp.displayname} placeholder='Room name' onChange={this.handleInputChange}/>
              <label>viewpermission</label>              
              <input className="form-text-input" name="viewpermission" value={this.state.resp.viewpermission} placeholder='viewpermission' onChange={this.handleInputChange}/>
              <label>bookpermission</label>              
              <input className="form-text-input" name="bookpermission" value={this.state.resp.bookpermission} placeholder='bookpermission' onChange={this.handleInputChange}/>
              <label>capacity</label>              
              <input className="form-text-input" name="capacity" value={this.state.resp.capacity} placeholder='0' onChange={this.handleInputChange}/>
            </Form.Field>
            <Button type='submit' onClick={() =>this.updateResource()}>Update</Button>
        </Form>
        <p>{JSON.stringify(serverresp)}</p>
      </div>
      }
        return (
          <div className="ManageRes">
          <h2>Room Details</h2>
          {output}
          </div>
        )
      
    }
  }

  class ManageRes extends React.Component {
    render() {

      let listItems = [];
      let listArea = null;
      if (this.state.resp != null) {
        listArea = <List divided relaxed className="item-list resource-list"> {listItems} </List>
        this.state.resp.forEach(function(element){
          listItems.push(
            <List.Item key={element.id}>
              <List.Icon name='user outline' size='large' verticalAlign='middle' />
              <List.Content>
                <List.Header>
                </List.Header>
                <List.Description>
                    <span>Name: {element.displayname}</span><br/>

                    <Button className="editUserBtn" as = {Link} to={"/adminview/users/detail/"+element.UserID}>Edit this room</Button>
                  
                </List.Description>
              </List.Content>
            </List.Item>
          );
        });
      } else {
        listArea = <p>Please search or browse rooms.</p>
      }
      
        return (
          <div className="ManageRes">
            <p>Manage Resources</p>

            <p>Add resource type</p>
            <Form action={serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/createNewresourceType'} method="post">
            <Form.Field>
              <label>resource type name</label>
              <input placeholder='typename' />
            </Form.Field>
            <Form.Field>
              <label>Last Name</label>
              <input placeholder='Last Name' />
            </Form.Field>
            <Form.Field>
              <Checkbox label='I agree to the Terms and Conditions' />
            </Form.Field>
            <Button type='submit'>Submit</Button>
          </Form>
          


            <p>Remove room1</p>
            <p>Edit room</p>
            {listArea}

          </div>
          )      
    }
  }

  class ManageUserMain extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        userid: '',
        email: '',
        resp: null,
      };
  
      this.handleInputChange = this.handleInputChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange(event) {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;

      this.setState({
        [name]: value
      });
    }
  
    handleSubmit(event) {
      alert('A name was submitted: ' + this.state.value);
      event.preventDefault();
    }

    findUser(offset) {
      var resp = null;
      var self = this;
      axios.post(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/finduser', {
          userid : self.state.userid,
          email : self.state.email,
          offset: offset,
          pageSize:20
        }, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          resp: response.data.Data
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 
    }

    listUser(offset) {
      var resp = null;
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/listusers?pageSize='+pageSize+'&offset=' + (offset*pageSize), {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          resp: response.data.Data
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 
    }


    render() {
      let listItems = [];
      let listArea = null;
      if (this.state.resp != null) {
        listArea = <List divided relaxed className="item-list user-list"> {listItems} </List>
        this.state.resp.forEach(function(element){
          let status = "Disabled"
          let statLabelColor = ""
          if (!element.Disabled) {
            status = "Enabled"
            statLabelColor = "blue"
          }
          listItems.push(
            <List.Item key={element.UserID}>
              <List.Icon name='user outline' size='large' verticalAlign='middle' />
              <List.Content>
                <List.Header>
                </List.Header>
                <List.Description>
                    <span>UserID: {element.UserID}</span><br/>
                    <span>Status: {status}</span><br/>
                    <span>Email: {element.Email}</span><br/>
                    <Button className="editUserBtn" as = {Link} to={"/adminview/users/detail/"+element.UserID}>Edit this user</Button>
                  
                </List.Description>
              </List.Content>
            </List.Item>
          );
        });
      } else {
        listArea = <p>Please search or browse users.</p>
      }
        return (
          <div className="ManageUserMain">
            <div className="searchField">
              <Form>
                <Form.Field>
                  {/*<label>User ID</label>*/}
                  <input className="form-text-input" name="userid" value={this.state.userid} placeholder='User ID' onChange={this.handleInputChange}/>
                </Form.Field>
                <Form.Field>
                  {/*<label>Email</label>*/}
                  <input className="form-text-input" name="email" value={this.state.email} placeholder='Email' onChange={this.handleInputChange}/>
                </Form.Field>            
                <Button type='submit' onClick={() =>this.findUser(0)}>Search</Button>
                <Button onClick={() =>this.listUser(0)}>Browse Users</Button>
              </Form>
            </div>
            
            {listArea}
          </div>
          )
      
    }
  }

  class ManageUser extends React.Component {
    render() {
      
        return (
          <div className="ManageUserView">
            <h4>Manage Users</h4>
            <Route exact path={'/adminview/users'} component={ManageUserMain}/>
            <Route path={'/adminview/users/detail/:uid'} component={UserDetails}/>
          </div>
          )
      
    }
  }

  class ManageRoom extends React.Component {
    render() {      
        return (
          <div className="ManageRoomView">
            <h4>Manage Rooms</h4>
            <Route exact path={'/adminview'} component={ManageRoomMain}/>
            <Route exact path={'/adminview/rooms'} component={ManageRoomMain}/>
            <Route path={'/adminview/rooms/detail/:rid'} component={RoomDetails}/>
          </div>
          )
      
    }
  }

  class ManageBooking extends React.Component {
    render() {      
        return (
          <div className="ManageBookingView">
            <h4>Manage Booking</h4>
            <Route exact path={'/adminview/booking'} component={ManageBookingMain}/>
            <Route path={'/adminview/bookings/detail/:bid'} component={BookingDetails}/>
          </div>
          )
      
    }
  }
  class ManageBookingMain extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        resp: null,
      };
      this.handleInputChange = this.handleInputChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange(event) {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;

      this.setState({
        [name]: value
      });
    }
  
    handleSubmit(event) {
      alert('A name was submitted: ' + this.state.value);
      event.preventDefault();
    }

    listBooking(today) {
      var resp = null;
      var self = this;
      var ep = "getbookinglistadm";
      if (today)
        ep = "getbookinglisttodayadm";

      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/'+ep, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        self.setState({
          resp: response.data.Data
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 

    }

    delBooking(id) {
      console.log("enter delbooking", id);
      var resp = null;
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/cancelbookingadm?bookingId=' + id, {withCredentials: true})
      .then(function (response) {
        if (response.data.Data == "OK"){
          var temp = [];
          for(var i = 0; i < self.state.resp.length; i++) {
            if(self.state.resp[i].id !== id) {
              temp.push(self.state.resp[i]);
            }
          }
          self.setState({
            resp: temp
          });
        }
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      }); 
    
    }


    render() {      
      let listItems = [];
      let listArea = null;
      if (this.state.resp != null && Array.isArray(this.state.resp)) {
        listArea = <List divided relaxed className="item-list user-list"> {listItems} </List>
        var self = this;
        this.state.resp.forEach(function(element){
          listItems.push(
            <List.Item key={element.id}>
              <List.Icon name='user outline' size='large' verticalAlign='middle' />
              <List.Content>
                <List.Header>
                </List.Header>
                <List.Description>
                    <span>Booked for: {element.bookedforuser}</span><br/>                    
                    <span>Room: {element.displayname}</span><br/>
                    <span>Start: {moment(element.bookstart).format("YYYY-MM-DD hh:mm a")}</span><br/>
                    <span>End: {moment(element.bookend).format("YYYY-MM-DD hh:mm a")}</span><br/>
                    <Button className="editUserBtn" as = {Link} to={"/adminview/bookings/detail/"+element.id}>Edit this booking</Button>
                    <Button className="delBookingBtn" onClick={self.delBooking.bind(self, element.id)}>Remove this booking</Button>
                  
                </List.Description>
              </List.Content>
            </List.Item>
          );
        });
      } else {
        listArea = <p>Please search or browse users.</p>
      }
        return (
          <div className="ManageBookingMainView">
            <h5>Manage Booking Main</h5>
            <Button onClick={() =>this.listBooking(true)}>Browse Bookings</Button>
            <div className="searchField">
              <Form>
                <Form.Field>
                  {/*<label>User ID</label>*/}
                  <input className="form-text-input" name="userid" value={this.state.userid} placeholder='User ID' onChange={this.handleInputChange}/>
                </Form.Field>
                <Form.Field>
                  {/*<label>Email</label>*/}
                  <input className="form-text-input" name="email" value={this.state.email} placeholder='Email' onChange={this.handleInputChange}/>
                </Form.Field>            
                <Button type='submit' onClick={() =>this.findUser(0)}>Search</Button>
                <Button onClick={() =>this.listUser(0)}>Browse Users</Button>
              </Form>
            </div>
            
            {listArea}
          </div>
          )
      
    }
  }

  class BookingDetails extends React.Component {
    render() {      
        return (
          <div className="BookingDetailsView">
            <h5>Manage Booking Details</h5>

          </div>
          )
      
    }
  }

  class UserDetails extends React.Component {
    constructor(props) {
      super(props);  
      this.state = { 
        resp: null,
        lastUID: ""
      };
      this.handleInputChange = this.handleInputChange.bind(this);
      //this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange = (event) => {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;
      let updatedresp = Object.assign({}, this.state.resp);    //creating copy of object
      updatedresp[name] = value;                        //updating value

      console.log("onchange fired", event.target);
      console.log("onchange fired", event.target.checked);
      console.log("onchange fired", event.target.value);

      this.setState({
        resp: updatedresp
      });
    }

    handleToggleChange = (event) => {
      const target = event.target;
      //const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.textContent;
      let updatedresp = Object.assign({}, this.state.resp);    //creating copy of object
      updatedresp[name] = !updatedresp[name];                        //updating value

      console.log("onchange fired", event.target);
      console.log("onchange fired", event.target.checked);


      this.setState({
        resp: updatedresp
      });
    }

    handleRoleInputChange = (event) => {
      const target = event.target;
      //const value = target.checked;
      const name = target.name;

      console.log("handleRoleInputChange fired", event.target);
      console.log("handleRoleInputChange fired", event.target.checked);
      console.log("handleRoleInputChange fired", event.target.value);

      let epName = (event.target.checked)? "assignroletouser" : "removerolefromuser";
      var self = this;
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/'+epName+'?username='+this.props.match.params.uid+'&rolename=' + name, {withCredentials: true})
      .then(function (response) {
        console.log(response);
        if (response.data.Data === "OK"){
          let updatedresp = Object.assign({}, self.state.resp);    //creating copy of object
          updatedresp.Roles[name] = !updatedresp.Roles[name];                        //updating value
          self.setState({
            resp: updatedresp
          });
        }
      })
      .catch(function (error) {
        console.log(error);
      });       
    }
  
    handleCancel = () => {      
      this.getData();
    }

    getData = () => {
      var resp = null;
      let uID = this.props.match.params.uid
      var self = this;
      if (uID) {
        axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/getuserdetails?uid='+uID, {withCredentials: true})
      .then(function (response) {
        console.log("response", response);
        self.setState({
          resp: response.data.Data,
          lastUID:uID
        });
        console.log("this.state.resp", self.state.resp);
        resp = response;

      })
      .catch(function (error) {
        console.log(error);
      });
      }      
    }

    toggleUser = () => {
      console.log("toggleUser");
      let endpoint = "disableuser"
      if (this.state.resp&&this.state.resp.Enabled){
        endpoint = "enableuser"
      } 
      axios.get(serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/'+ endpoint +'?username='+this.props.match.params.uid, {withCredentials: true}).then(function (response) {
        console.log("response", response);        
      })
      .catch(function (error) {
        console.log(error);
      });

    }

    componentDidMount() {
      this.getData()
    }

    componentDidUpdate(){
      if (this.props.match.params.uid != this.state.lastUID) {
        this.getData()
      }

    }

    render() {
      let uID = this.props.match.params.uid;
      let serverresp = this.state.resp;
      let output = <div></div>;
      let roleuis = [];

      let statusOutput = <div>
        <p>Status: Disabled</p>
        <Button primary onClick={() =>this.toggleUser()}>Enable this user</Button>
      </div>

      if (this.state.resp&&this.state.resp.Enabled){
        statusOutput = <div>
        <p>Status: Enabled</p>
        <Button secondary onClick={() =>this.toggleUser()}>Disable this user</Button>
      </div>
      }
      
      if (serverresp!=null){
        for (var rol in serverresp.Roles){
          if (!serverresp.Roles.hasOwnProperty(rol)) continue;
          roleuis.push(
            <Form.Field key = {rol}>
            <label> {rol} &nbsp;
              <input name={rol} type="checkbox" checked={this.state.resp.Roles[rol]} onChange={this.handleRoleInputChange} />
            </label>            
            </Form.Field>
          );
        }
        output= <div className="UserDetails">
        <h2>UserDetails</h2>   
        <p>UID: {uID}</p>
        <p>Email: {serverresp.Email}</p>
        <p>Created: {serverresp.Created}</p>
        <p>Lastlogin: {serverresp.Lastlogin}</p>
        {statusOutput}

        <br/>
        <Form>
        <h3>Roles</h3>
        {roleuis}
        {/*<Button type='submit' onClick={() =>this.submitRoleChanges()}>Save Role Change</Button>*/}
        </Form>
        <p>{JSON.stringify(serverresp)}</p>
      </div>
      }
        return (output);
      
    }
  }

  class ShowStats extends React.Component {
    render() {
      
        return (
          <div className="ShowStats">
            <p>ShowStats</p>

          </div>
          )
      
    }
  }

  class LoginView extends React.Component {
    render() {
      
        return (
          <div className="LoginView">
            <p>LoginView</p>
            <form action={serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/login'} method="post">
            Username:<input type="text" name="username"/>
            Password:<input type="password" name="password"/>
            <input type="submit" value="Login"/>
        </form>

          </div>
          )
      
    }
  }

  class SignupView extends React.Component {
    render() {
      
        return (
          <div className="SignupView">
            <p>SignupView</p>
            <form action={serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/createuserbysignup'} method="post">
            Username:<input type="text" name="username"/>
            Email:<input type="text" name="email"/>
            Password:<input type="password" name="password"/>
            Repeat Password:<input type="password"/>
            <input type="submit" value="Create account"/>
        </form>

          </div>
          )
      
    }
  }
  export default App;
