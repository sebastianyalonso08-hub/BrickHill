const express=require("express");
const session=require("express-session");
const fs=require("fs");
const path=require("path");
const http=require("http");
const crypto=require("crypto");
const {WebSocketServer}=require("ws");

const app=express();
const server=http.createServer(app);
const PORT=Number(process.env.PORT||3000);
const DATA=path.join(__dirname,"data");
fs.mkdirSync(DATA,{recursive:true});
const dbFile=path.join(DATA,"site.json");
if(!fs.existsSync(dbFile))fs.writeFileSync(dbFile,JSON.stringify({users:{},games:[]},null,2));
function db(){return JSON.parse(fs.readFileSync(dbFile,"utf8"))}
function save(x){fs.writeFileSync(dbFile,JSON.stringify(x,null,2))}
app.use(express.json());app.use(express.urlencoded({extended:true}));
app.use(session({secret:process.env.SESSION_SECRET||"brick-hill-v2-dev",resave:false,saveUninitialized:false,cookie:{httpOnly:true,sameSite:"lax",maxAge:7*24*60*60*1000}}));
app.use(express.static(path.join(__dirname,"public")));
app.get("/api/client/package",(req,res)=>{
  const pkg=path.join(__dirname,"public","BrickHillClient.zip");
  if(!fs.existsSync(pkg))return res.status(404).send("Client package unavailable.");
  res.download(pkg,"BrickHillClient.zip");
});

const tickets=new Map();
const connections=new Map();
const rooms=new Map();
function token(){return crypto.randomBytes(32).toString("base64url")}
function newTicket(user,game){const t=token();tickets.set(t,{user,game,expires:Date.now()+60000,used:false});return t}
function consumeTicket(t){const x=tickets.get(t);if(!x||x.used||x.expires<Date.now())return null;x.used=true;tickets.delete(t);return x}
function newConnection(t){const c=token();connections.set(c,{...t,expires:Date.now()+120000});return c}
function getConnection(c){const x=connections.get(c);if(!x||x.expires<Date.now())return null;return x}

app.get("/api/me",(req,res)=>{const d=db();res.json({user:req.session.user?d.users[req.session.user]||null:null})});
app.post("/api/register",(req,res)=>{const username=String(req.body.username||"").trim();if(!/^[A-Za-z0-9_]{3,20}$/.test(username))return res.status(400).json({error:"Username must be 3-20 characters."});const d=db(),key=username.toLowerCase();if(d.users[key])return res.status(400).json({error:"Username is already taken."});d.users[key]={id:"USR-"+String(Object.keys(d.users).length+1).padStart(5,"0"),username,joined:Date.now(),bricks:100};save(d);req.session.user=key;res.json({ok:true,user:d.users[key]})});
app.post("/api/login",(req,res)=>{const key=String(req.body.username||"").trim().toLowerCase(),d=db();if(!d.users[key])return res.status(401).json({error:"Invalid username."});req.session.user=key;res.json({ok:true,user:d.users[key]})});
app.post("/api/logout",(req,res)=>req.session.destroy(()=>res.json({ok:true})));
app.get("/api/games",(req,res)=>{const d=db();const online={};for(const [id,room] of rooms)online[id]=room.size;res.json({games:d.games.map(g=>({...g,players:online[g.id]||0}))})});
app.post("/api/games",(req,res)=>{if(!req.session.user)return res.status(401).json({error:"Log in first."});const d=db(),u=d.users[req.session.user],name=String(req.body.name||"").trim();if(!name||name.length>50)return res.status(400).json({error:"Enter a game name."});const game={id:String(Date.now()),name,creator:u.username,players:0,visits:0,description:String(req.body.description||"").slice(0,180),featured:false};d.games.unshift(game);save(d);res.json({ok:true,game})});

app.post("/api/client/launch",(req,res)=>{if(!req.session.user)return res.status(401).json({error:"Log in first."});const gameId=String(req.body.gameId||"");const d=db();const game=d.games.find(g=>g.id===gameId);if(!game)return res.status(404).json({error:"Game not found."});const ticket=newTicket(d.users[req.session.user],gameId);res.json({ok:true,scheme:`brickhill://play?game=${encodeURIComponent(gameId)}&ticket=${encodeURIComponent(ticket)}`,game:{id:game.id,name:game.name}})});
app.post("/api/client/redeem",(req,res)=>{const t=consumeTicket(String(req.body.ticket||""));if(!t)return res.status(401).json({error:"Invalid or expired launch ticket."});const connection=newConnection(t);res.json({ok:true,user:t.user,gameId:t.game,ws:`wss://${req.get("host")}/ws/legacy?game=${encodeURIComponent(t.game)}&connection=${encodeURIComponent(connection)}`})});
app.get("/api/client/installer",(req,res)=>{
  const installer=path.join(__dirname,"install-brickhill.ps1");
  if(!fs.existsSync(installer))return res.status(404).send("Installer unavailable.");
  res.download(installer,"install-brickhill.ps1");
});
app.get("/health",(req,res)=>res.json({ok:true,version:"2.2.0",websocket:true,legacyBridge:true}));

const wss=new WebSocketServer({server,path:"/ws/legacy"});
wss.on("connection",(ws,req)=>{
  const u=new URL(req.url,"http://localhost");
  const game=u.searchParams.get("game")||"default";
  const connection=getConnection(u.searchParams.get("connection")||"");
  if(!connection || connection.game!==game){ws.close(1008,"invalid connection");return}
  if(!rooms.has(game))rooms.set(game,new Set());
  const room=rooms.get(game);room.add(ws);ws._bhGame=game;ws._bhUser=connection.user;
  const leave=()=>{room.delete(ws);if(!room.size)rooms.delete(game)};
  ws.on("message",data=>{for(const peer of room){if(peer!==ws&&peer.readyState===1)peer.send(data,{binary:true})}});
  ws.on("close",leave);ws.on("error",leave);
});
setInterval(()=>{const now=Date.now();for(const [k,v] of connections)if(v.expires<now)connections.delete(k);for(const [k,v] of tickets)if(v.expires<now)tickets.delete(k)},30000);
server.listen(PORT,"0.0.0.0",()=>console.log(`Brick Hill Website V2.1 running on ${PORT}`));
