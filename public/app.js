let me=null;
const pages=[...document.querySelectorAll(".page")];
function go(id){pages.forEach(p=>p.classList.toggle("active",p.id===id));window.scrollTo(0,0);if(id==="games")loadGames();if(id==="home")loadHome();if(id==="profile")loadProfile()}
document.addEventListener("click",e=>{let b=e.target.closest("[data-page]");if(b){e.preventDefault();go(b.dataset.page)}});
const tutorialButton=document.getElementById("openTutorial");
if(tutorialButton)tutorialButton.onclick=()=>go("tutorial");
async function api(url,opt){let r=await fetch(url,opt),d=await r.json();if(!r.ok)throw Error(d.error||"Request failed");return d}
function esc(x){return String(x).replace(/[&<>"']/g,c=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]))}
function renderAccount(){document.getElementById("logoutTop").style.display=me?"block":"none";document.getElementById("welcomeName").textContent=me?me.username:"Guest";document.getElementById("homeUser").textContent=me?me.username:"Guest";document.getElementById("homeBricks").textContent=me?me.bricks:0;document.getElementById("miniUser").innerHTML=me?`<b>${esc(me.username)}</b><div class="coins">🧱 ${me.bricks}　★ 0</div>`:"<b>Welcome!</b><div class='coins'>🧱 0　★ 0</div>"}
async function loadHome(){try{let d=await api("/api/games");document.getElementById("featured").innerHTML=d.games.slice(0,3).map(g=>`<div class="game-card"><div class="thumb"></div><div><h3>${esc(g.name)}</h3><p>${esc(g.description)}</p><p>by ${esc(g.creator)} · ${g.players||0} playing</p></div><button class="playbtn" data-play="${esc(g.id)}">Play</button></div>`).join("")}catch{}}
async function loadGames(){try{let d=await api("/api/games");document.getElementById("gamesList").innerHTML=d.games.map(g=>`<div class="game-card"><div class="thumb"></div><div><h3>${esc(g.name)}</h3><p>${esc(g.description)}</p><p>by ${esc(g.creator)} · ${g.players||0} playing</p></div><button class="playbtn" data-play="${esc(g.id)}">Play</button></div>`).join("")}catch{}}
async function loadProfile(){let b=document.getElementById("profileBox");if(!me){b.innerHTML=`Please log in. <button class="tinybtn" data-page="login">Login</button>`;return}document.getElementById("profileName").textContent=me.username;b.innerHTML=`<b>${esc(me.username)}</b><br>ID: ${esc(me.id)}<br>Bricks: ${me.bricks}<br><br><button class="tinybtn">Send Message</button> <button class="tinybtn">Add Friend</button>`}
document.getElementById("register").onclick=async()=>{try{await api("/api/register",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({username:document.getElementById("registerName").value})});await refresh();go("home")}catch(e){document.getElementById("registerError").textContent=e.message}}
document.getElementById("login").onclick=async()=>{try{await api("/api/login",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({username:document.getElementById("loginName").value})});await refresh();go("home")}catch(e){document.getElementById("loginError").textContent=e.message}}
document.getElementById("logoutTop").onclick=async()=>{await api("/api/logout",{method:"POST"});me=null;renderAccount();go("home")}
document.getElementById("publish").onclick=async()=>{try{await api("/api/games",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({name:document.getElementById("gameName").value,description:document.getElementById("gameDesc").value})});document.getElementById("gameName").value="";document.getElementById("gameDesc").value="";loadGames()}catch(e){document.getElementById("gameError").textContent=e.message}}
async function launchGame(gameId){
  if(!me){go("login");return}
  try{
    const d=await api("/api/client/launch",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({gameId})});
    go("client");
    const status=document.getElementById("clientStatus");
    // The custom protocol only works after the Windows launcher has been installed.
    if(status) status.innerHTML=`Starting <b>${esc(d.game.name)}</b>… If nothing opens, use <a href="/api/client/installer" target="_blank" rel="noopener">Install Brick Hill Client</a>, then click Play again.`;
    const link=document.createElement("a");
    link.href=d.scheme;
    link.setAttribute("aria-hidden","true");
    link.style.display="none";
    document.body.appendChild(link);
    link.click();
    setTimeout(()=>link.remove(),1000);
  }catch(e){go("client");const status=document.getElementById("clientStatus");if(status)status.innerHTML=`<b>Could not start the client.</b><br>${esc(e.message)}<br><br><a href="/api/client/installer" target="_blank" rel="noopener">Install Brick Hill Client</a>`}
}
document.addEventListener("click",e=>{let p=e.target.closest("[data-play]");if(p)launchGame(p.dataset.play)});
async function refresh(){let d=await api("/api/me");me=d.user;renderAccount();loadHome();loadGames();loadProfile()}
refresh();
