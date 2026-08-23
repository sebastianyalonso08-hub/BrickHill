const express=require("express");
const session=require("express-session");
const fs=require("fs");
const path=require("path");

const app=express();
const PORT=Number(process.env.PORT||3000);
const DATA=path.join(__dirname,"data");
fs.mkdirSync(DATA,{recursive:true});

const dbFile=path.join(DATA,"site.json");
if(!fs.existsSync(dbFile)){
  fs.writeFileSync(dbFile,JSON.stringify({
    users:{},
    games:[
      {id:"1001",name:"Brick Hill Classic",creator:"Brick Hill",players:0,visits:0,description:"A classic Brick Hill experience.",featured:true},
      {id:"1002",name:"Brick Obby",creator:"Brick Hill",players:0,visits:0,description:"Jump across a world of bricks.",featured:true},
      {id:"1003",name:"Build & Chill",creator:"Brick Hill",players:0,visits:0,description:"Build something and hang out.",featured:false}
    ]
  },null,2));
}
function db(){return JSON.parse(fs.readFileSync(dbFile,"utf8"))}
function save(x){fs.writeFileSync(dbFile,JSON.stringify(x,null,2))}

app.use(express.json());
app.use(express.urlencoded({extended:true}));
app.use(session({
  secret:process.env.SESSION_SECRET||"brick-hill-v1-dev",
  resave:false,saveUninitialized:false,
  cookie:{httpOnly:true,sameSite:"lax",maxAge:7*24*60*60*1000}
}));
app.use(express.static(path.join(__dirname,"public")));

app.get("/api/me",(req,res)=>{
  const d=db();
  res.json({user:req.session.user?d.users[req.session.user]||null:null});
});

app.post("/api/register",(req,res)=>{
  const username=String(req.body.username||"").trim();
  if(!/^[A-Za-z0-9_]{3,20}$/.test(username))
    return res.status(400).json({error:"Username must be 3-20 characters."});
  const d=db(), key=username.toLowerCase();
  if(d.users[key]) return res.status(400).json({error:"Username is already taken."});
  d.users[key]={id:"USR-"+String(Object.keys(d.users).length+1).padStart(5,"0"),username,joined:Date.now(),bricks:100};
  save(d); req.session.user=key;
  res.json({ok:true,user:d.users[key]});
});

app.post("/api/login",(req,res)=>{
  const key=String(req.body.username||"").trim().toLowerCase();
  const d=db();
  if(!d.users[key]) return res.status(401).json({error:"Invalid username."});
  req.session.user=key; res.json({ok:true,user:d.users[key]});
});

app.post("/api/logout",(req,res)=>req.session.destroy(()=>res.json({ok:true})));

app.get("/api/games",(req,res)=>res.json({games:db().games}));

app.post("/api/games",(req,res)=>{
  if(!req.session.user) return res.status(401).json({error:"Log in first."});
  const d=db(), u=d.users[req.session.user];
  const name=String(req.body.name||"").trim();
  if(!name||name.length>50) return res.status(400).json({error:"Enter a game name."});
  const game={id:String(Date.now()),name,creator:u.username,players:0,visits:0,description:String(req.body.description||"").slice(0,180),featured:false};
  d.games.unshift(game); save(d); res.json({ok:true,game});
});

app.get("/health",(req,res)=>res.json({ok:true,version:"1.0.0"}));

app.listen(PORT,"0.0.0.0",()=>console.log(`Brick Hill Website V1 running at http://localhost:${PORT}`));