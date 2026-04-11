package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// serves the browser dashboard at GET /
func (h *Handler) UI(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	c.Writer.WriteString(uiHTML)
}

const uiHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Olu SMS Mock</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#f9fafb;color:#111827;font-size:14px}
.nav{background:#111827;color:#fff;padding:0 20px;height:52px;display:flex;align-items:center;justify-content:space-between}
.nav-brand{font-size:15px;font-weight:600;display:flex;align-items:center;gap:8px}
.pulse{width:8px;height:8px;border-radius:50%;background:#22c55e;animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.3}}
.nav-badge{font-size:11px;background:#22c55e22;color:#22c55e;border:1px solid #22c55e44;padding:2px 10px;border-radius:20px}
.nav-info{font-size:12px;color:#6b7280}
.container{max-width:860px;margin:0 auto;padding:20px}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:10px;margin-bottom:18px}
.stat{background:#fff;border:1px solid #e5e7eb;border-radius:10px;padding:12px 16px}
.stat-n{font-size:22px;font-weight:600}
.stat-n.green{color:#008751}
.stat-n.blue{color:#1d4ed8}
.stat-n.amber{color:#92400e}
.stat-n.red{color:#991b1b}
.stat-l{font-size:11px;color:#6b7280;margin-top:2px}
.toolbar{display:flex;gap:10px;margin-bottom:14px;align-items:center}
.search{flex:1;height:36px;border:1px solid #d1d5db;border-radius:8px;padding:0 12px;font-size:13px;outline:none;background:#fff}
.search:focus{border-color:#008751;box-shadow:0 0 0 3px rgba(0,135,81,.08)}
.btn{height:36px;padding:0 16px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;border:1px solid #d1d5db;background:#fff;color:#374151}
.btn:hover{background:#f3f4f6}
.btn.danger{border-color:#fca5a5;color:#991b1b}
.btn.danger:hover{background:#fef2f2}
.btn.primary{background:#008751;color:#fff;border-color:#008751}
.btn.primary:hover{background:#007347}
.msgs{display:flex;flex-direction:column;gap:8px}
.msg{background:#fff;border:1px solid #e5e7eb;border-radius:10px;padding:14px 16px;border-left-width:4px}
.msg.inbound{border-left-color:#7c3aed}
.msg.otp{border-left-color:#008751}
.msg.confirm{border-left-color:#1d4ed8}
.msg.reject{border-left-color:#dc2626}
.msg.generic{border-left-color:#d1d5db}
.msg-header{display:flex;align-items:center;justify-content:space-between;margin-bottom:8px}
.msg-phone{font-size:13px;font-weight:600;font-family:monospace;color:#111827}
.msg-time{font-size:11px;color:#9ca3af}
.msg-from{font-size:11px;color:#9ca3af;margin-bottom:6px}
.msg-body{font-size:13px;color:#374151;line-height:1.5;background:#f9fafb;border-radius:6px;padding:8px 10px}
.otp-code{font-family:monospace;font-size:16px;font-weight:700;color:#008751;background:#ecfdf5;padding:2px 8px;border-radius:4px;border:1px solid #6ee7b7}
.tag{font-size:10px;font-weight:600;padding:2px 8px;border-radius:20px;margin-left:8px}
.tag.inbound{background:#f3e8ff;color:#6d28d9}
.tag.otp{background:#ecfdf5;color:#065f46}
.tag.confirm{background:#eff6ff;color:#1e40af}
.tag.reject{background:#fef2f2;color:#991b1b}
.tag.generic{background:#f3f4f6;color:#6b7280}
.tag.source{border:1px solid transparent;text-transform:uppercase}
.tag.source.web{background:#e0f2fe;color:#075985;border-color:#bae6fd}
.tag.source.sms{background:#fffbeb;color:#92400e;border-color:#fde68a}
.empty{text-align:center;padding:48px;color:#9ca3af;background:#fff;border:1px solid #e5e7eb;border-radius:10px}
.copy-btn{font-size:11px;padding:1px 8px;border:1px solid #d1d5db;border-radius:4px;background:#fff;cursor:pointer;color:#6b7280;margin-left:8px}
.copy-btn:hover{background:#f3f4f6}
.copy-btn.copied{background:#ecfdf5;color:#065f46;border-color:#6ee7b7}
footer{text-align:center;padding:20px;font-size:11px;color:#9ca3af}
</style>
</head>
<body>
<nav class="nav">
  <div class="nav-brand"><div class="pulse"></div>Olu SMS mock</div>
  <div style="display:flex;align-items:center;gap:16px">
    <span class="nav-info" id="nav-host"></span>
    <span class="nav-badge">running</span>
  </div>
</nav>
 
<div class="container">
  <div class="stats">
    <div class="stat"><div class="stat-n green" id="s-total">0</div><div class="stat-l">Total sent</div></div>
    <div class="stat"><div class="stat-n blue" id="s-otp">0</div><div class="stat-l">OTPs</div></div>
    <div class="stat"><div class="stat-n" id="s-confirm">0</div><div class="stat-l">Confirmations</div></div>
    <div class="stat"><div class="stat-n red" id="s-reject">0</div><div class="stat-l">Rejections</div></div>
  </div>
 
  <div class="toolbar">
    <input class="search" id="search" type="text" placeholder="Filter by phone number..." oninput="render()">
    <button class="btn primary" onclick="load()">Refresh</button>
    <button class="btn danger" onclick="clearAll()">Clear all</button>
  </div>
 
  <div class="msgs" id="msgs">
    <div class="empty">No messages yet. Send an OTP request to see it here.</div>
  </div>
</div>
<footer>GET /otp/:phone &nbsp;&middot;&nbsp; GET /messages/:phone &nbsp;&middot;&nbsp; GET /messages/:phone/latest &nbsp;&middot;&nbsp; DELETE /messages</footer>
 
<script>
let all = [];
 
async function load(){
  try{
    const r = await fetch('/messages');
    const d = await r.json();
    all = d.messages || [];
    const stats = d.stats || {};
    document.getElementById('s-total').textContent = stats.total || 0;
    document.getElementById('s-otp').textContent = stats.otp || 0;
    document.getElementById('s-confirm').textContent = stats.confirm || 0;
    document.getElementById('s-reject').textContent = stats.reject || 0;
    render();
  }catch(e){console.error(e)}
}
 
function render(){
  const q = document.getElementById('search').value.toLowerCase();
  const filtered = q ? all.filter(m => m.to.includes(q) || m.from.includes(q) || m.body.toLowerCase().includes(q)) : all;
  const el = document.getElementById('msgs');
  if(!filtered.length){
    el.innerHTML = '<div class="empty">No messages yet. Send an OTP request to see it here.</div>';
    return;
  }
  el.innerHTML = filtered.map(m => {
    const t = new Date(m.sent_at).toLocaleTimeString('en-NG',{hour:'2-digit',minute:'2-digit',second:'2-digit'});
    const body = m.otp_code
      ? highlightOTP(m.body, m.otp_code)
      : escHtml(m.body);
    return '<div class="msg '+m.channel+'">'+
      '<div class="msg-header">'+
        '<span class="msg-phone">'+escHtml(m.to)+'<span class="tag '+m.channel+'">'+m.channel+'</span>'+sourceTag(m.source)+'</span>'+
        '<span class="msg-time">'+t+'</span>'+
      '</div>'+
      '<div class="msg-from">From: '+escHtml(m.from)+'</div>'+
      '<div class="msg-body">'+body+'</div>'+
    '</div>';
  }).join('');
}
 
function escHtml(s){
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

function sourceTag(source){
  if(source !== 'web' && source !== 'sms') return '';
  return '<span class="tag source '+source+'">'+source+'</span>';
}

function highlightOTP(body, otp){
  const escapedBody = escHtml(body);
  const escapedOTP = escHtml(otp);
  const otpMarkup = '<span class="otp-code">'+escapedOTP+'<button class="copy-btn" onclick="copyOTP('+JSON.stringify(otp)+',this)">copy</button></span>';
  return escapedBody.replace(escapedOTP, otpMarkup);
}
 
async function clearAll(){
  if(!confirm('Clear all messages?')) return;
  await fetch('/messages',{method:'DELETE'});
  all = [];
  ['total','otp','confirm','reject'].forEach(k => document.getElementById('s-'+k).textContent='0');
  render();
}
 
function copyOTP(code, btn){
  navigator.clipboard.writeText(code).then(()=>{
    btn.textContent='copied!';
    btn.classList.add('copied');
    setTimeout(()=>{btn.textContent='copy';btn.classList.remove('copied');},2000);
  });
}
 
document.getElementById('nav-host').textContent = window.location.host + ' \u2014 dev only';

load();
setInterval(load, 2000);
</script>
</body>
</html>`
