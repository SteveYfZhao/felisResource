import 'rc-calendar/assets/index.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Calendar from 'rc-calendar';
import logo from './logo.svg';
import './App.css';
import axios from 'axios';
import { BrowserRouter, Link, Route } from 'react-router-dom'
import { Button, Grid, Icon, Label, Menu, Table, Input, Divider, Select, Checkbox, Sidebar, Segment, Header, Image, Form, List } from 'semantic-ui-react';

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
          <header className="App-header">
            {/*<img src={logo} className="App-logo" alt="logo" />*/}
            <h1 className="App-title"><a href="/">Felis Resource Management System</a></h1>
          </header>
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
    render() {
      return (
        <div className="MainView">
          
          
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
          {/*<p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
         </p>*/}

          <Grid celled>
            <Grid.Row>
              <Grid.Column width={3}>
                <p>Select date</p>
                <Calendar />
                <div>
                  <Button>Only show rooms available</Button>
                  <p>for at least <input placeholder='1' size = '2'/> hour</p>
                  <div>from <Input type='text' size='mini' placeholder='12' action>
                  <input size = '2'/>
                  <Select compact options={halfhouroptions} defaultValue='00' />
                  <Select compact options={apoptions} defaultValue='PM' />
                </Input></div>

                  <Divider horizontal>Or</Divider>
                  <Button>Show all</Button>
                </div>
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
      var rNum = this.props.rNum;
      var colPerPage = this.props.colPerPage;
      var pageNum = this.props.pageNum;
      
      if (!rlist[0].SpanOverride){

        var key = rNum + "override";
        cells.push(<Table.Cell key={key} rowSpan = {rlist[0].rowSpan} selectable = {rlist[0].selectable}>{rlist[0].value}</Table.Cell>);
      }
      for (var i=0; i < this.props.colPerPage; i++) {
        var cellData = rlist[i + 1 + pageNum * colPerPage];
          if (cellData == null || cellData.SpanOverride){
            continue;
          } else {
            //console.log("push cell", cells);
            var key = rNum + "c" + i
            cells.push(<Table.Cell key = {key} rowSpan = {cellData.rowSpan} selectable = {cellData.selectable}>{cellData.value}</Table.Cell>);
            
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
          row.push({key:h, rowSpan:2, selectable:false, SpanOverride:false, value:momentParser.format("hh:mm a")});
        } else if (h % 2 == 1) {
          row.push({key:h, rowSpan:1, selectable:false, SpanOverride:true, value:""});
        }        
        momentParser.add(30, "m");

        for (let i = 0; i < totalItem; i++){
          var cellData = {rowSpan:1, selectable:true, SpanOverride:false, value:""};          
          row.push(cellData);
        }
        //console.log("datarow", row);
        rMatrixJSX.push(<TRowSlot key={"hour"+h} id={"hour"+h} rlist={row} colPerPage = {colPerPage} pageNum = {pageNum}/>)
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
              <Route exact path={'/adminview'} component={ManageRooms}/>
              <Route path={'/adminview/rooms'} component={ManageRooms}/>
              <Route path={'/adminview/resources'} component={ManageRes}/>
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

  class ManageRooms extends React.Component {
    render() {
      
        return (
          <div className="ManageRooms">
            <p>Manage rooms</p>

            <p>Add room</p>
            <Form action={serverProtocol + "://" + window.location.hostname + ':' + serverPortNum +'/addresource'} method="post">
            <Form.Field>
              <label>Room name</label>
              <input placeholder='Room name' />
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


            <p>Remove room</p>
            <p>Edit room</p>

          </div>
          )      
    }
  }

  class ManageRes extends React.Component {
    render() {
      
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


            <p>Remove room</p>
            <p>Edit room</p>

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
        listArea = <List divided relaxed> {listItems} </List>
        this.state.resp.forEach(function(element){
          let status = "Disabled"
          let statLabelColor = ""
          if (!element.Disabled) {
            status = "Enabled"
            statLabelColor = "blue"
          }
          listItems.push(
            <List.Item>
            <List.Icon name='user outline' size='large' verticalAlign='middle' />
            <List.Content>
              <List.Header as='a'>UserID: {element.UserID } 
                <Label color={statLabelColor}>{status}</Label>
              </List.Header>
              <List.Description as='a'>Email: {element.Email} <Link to={"/adminview/users/detail/"+element.UserID}>Edit this user</Link></List.Description>
            </List.Content>
          </List.Item>);
        });
      } else {
        listArea = <p>Please search or browse users.</p>
      }
        return (
          <div className="ManageUserMain">
            <Form>
            <Form.Field>
              <label>User ID</label>
              <input name="userid" value={this.state.userid} placeholder='User ID' onChange={this.handleInputChange}/>
            </Form.Field>
            <Form.Field>
              <label>Email</label>
              <input name="email" value={this.state.email} placeholder='Email' onChange={this.handleInputChange}/>
            </Form.Field>            
            <Button type='submit' onClick={() =>this.findUser(0)}>Search</Button>
          </Form>
          <Button onClick={() =>this.listUser(0)}>Browse Users</Button>
          {listArea}
          </div>
          )
      
    }
  }

  class ManageUser extends React.Component {
    render() {
      
        return (
          <div className="ManageUserView">
            <p>Manage Users</p>
            <Route exact path={'/adminview/users'} component={ManageUserMain}/>
            <Route path={'/adminview/users/detail/:uid'} component={UserDetails}/>
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
      this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleInputChange(event) {
      const target = event.target;
      const value = target.type === 'checkbox' ? target.checked : target.value;
      const name = target.name;
      let updatedresp = Object.assign({}, this.state.resp);    //creating copy of object
      updatedresp[name] = value;                        //updating value

      this.setState({
        resp: updatedresp
      });
    }

    handleRoleInputChange(event) {
      const target = event.target;
      const value = target.checked;
      const name = target.name;

      let updatedresp = Object.assign({}, this.state.resp);    //creating copy of object
      updatedresp.Roles[name] = value;                        //updating value

      this.setState({
        resp: updatedresp
      });
    }
  
    handleSubmit(event) {
      alert('A name was submitted: ' + this.state.value);
      event.preventDefault();
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

      

      
      if (serverresp!=null){
        for (var rol in serverresp.Roles){
          if (!serverresp.Roles.hasOwnProperty(rol)) continue;
          roleuis.push(
            <Form.Field>
            <Checkbox label = {rol} name={rol} checked={this.state.resp.Roles[rol]} onChange={this.handleRoleInputChange}/>
            </Form.Field>
          );
        }
        output= <div className="UserDetails">
        <h2>UserDetails</h2>   
        <p>UID: {uID}</p>
        <p>Created: {serverresp.Created}</p>
        <p>Lastlogin: {serverresp.Lastlogin}</p>

        <Form>
        <Form.Field>
          
          <Checkbox label = "Disabled" toggle name="disabled" checked={this.state.resp.Disabled} onChange={this.handleInputChange}/>          
        </Form.Field>
        <Form.Field>
          <label>Email</label>
          <input name="email" value={this.state.email} placeholder='Email' onChange={this.handleInputChange}/>
        </Form.Field>
        <h3>Roles</h3>
        {roleuis}
        <Button type='submit' onClick={() =>this.findUser(0)}>Search</Button>
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
