(this.webpackJsonprtsp2flvweb=this.webpackJsonprtsp2flvweb||[]).push([[0],{163:function(e,t,a){"use strict";a.r(t);var n=a(0),o=a.n(n),r=a(12),c=a.n(r),i=a(87),l=a(16),s=a(114),d=a(14),u=a(209),j=a(165),b=a(229),h=a(233),p=a(232),m=a(228),f=a(230),O=a(243),x=a(231),g=a(112),v=a.n(g),w=a(111),y=a.n(w),k=a(225),C=a(224),S=a(226),N=a(227),E=a(110),I=a.n(E),A=a(105),R=a.n(A),T=a(216),B=a(245),P=a(215),F=a(213),L=a(212),M=a(214),D=a(2),W=Object(u.a)((function(e){return{root:{width:"100%"}}}));function V(e){var t=W(),a=o.a.useState(!1),r=Object(d.a)(a,2),c=r[0],i=r[1],l=o.a.useState(e.dialog.title),s=Object(d.a)(l,2),u=s[0],j=(s[1],o.a.useState(e.dialog.content)),b=Object(d.a)(j,2),h=b[0];b[1];Object(n.useImperativeHandle)(e.onRef,(function(){return{handleClickOpen:p}}));var p=function(){i(!0)},m=function(){i(!1)};return Object(D.jsx)("div",{className:t.root,children:Object(D.jsxs)(B.a,{open:c,onClose:m,"aria-labelledby":"form-dialog-title",children:[Object(D.jsx)(L.a,{id:"form-dialog-title",children:u}),Object(D.jsx)(F.a,{children:Object(D.jsx)(M.a,{children:h})}),Object(D.jsx)(P.a,{children:Object(D.jsx)(T.a,{onClick:m,color:"primary",children:"OK"})})]})})}var H=a(80),U=a(218),q=a(219),z=a(40),Y=a(220),J=a(67),_=a.n(J),K=a(217),$=a(48),G={serverURL:window.location.origin},Q=a(34),X=a.n(Q);X.a.defaults.withCredentials=!1,X.a.defaults.timeout=1e4,X.a.interceptors.request.use((function(e){var t=window.localStorage.getItem("token");return e.headers.token=t,e}),(function(e){return console.error(e),Promise.reject(e)})),X.a.interceptors.response.use((function(e){return e.data&&e.data.errcode&&401===parseInt(e.data.errcode)&&(window.location.hash="#/login"),e}),(function(e){return X.a.isCancel()||(e.response?(console.log(401===e.response.status),401===e.response.status?window.localStorage.getItem("token")&&window.localStorage.getItem("tokenExpired")&&"false"!==window.localStorage.getItem("tokenExpired")||(window.localStorage.setItem("tokenExpired","true"),window.location.hash="#/login"):500===e.response.status&&alert("server exception !")):e&&"error: timeout"===String(e).toLowerCase().substring(0,14)?alert("server timeout !"):alert("server error !")),Promise.reject(e)}));var Z=G.serverURL,ee=function(e,t){return X.a.post("".concat(Z).concat(e),t).then((function(e){return e.data}))},te="".concat(G.serverURL),ae=function(e){return ee("/system/login",e)},ne=function(e){return function(e,t){return X.a.get("".concat(Z).concat(e),{params:t}).then((function(e){return e.data}))}("/camera/list",e)},oe=function(e){return ee("/camera/edit",e)},re=function(e){return ee("/camera/delete/".concat(e.id),e)},ce=function(e){return ee("/camera/enabled",e)},ie=Object(u.a)((function(e){return{appBar:{position:"relative",backgroundColor:"#eebbaa"},title:{marginLeft:e.spacing(2),flex:1},videoContainer:{width:"90%",margin:"0 auto"}}})),le=o.a.forwardRef((function(e,t){return Object(D.jsx)(K.a,Object(H.a)({direction:"up",ref:t},e))}));function se(e){var t=ie(),a=o.a.useState(!1),r=Object(d.a)(a,2),c=r[0],i=r[1],l=o.a.useState(!0),s=Object(d.a)(l,2),u=s[0],j=s[1],b=null,h=0;Object(n.useImperativeHandle)(e.onRef,(function(){return{handleClickOpen:p}}));var p=function(){i(!0)},m=function(){i(!1)},f=function(t){var a={type:"flv"},n=te+"/live/permanent/"+e.row.code+"/"+e.row.playAuthCode+".flv";a.url=n,a.hasAudio=u,a.isLive=!0,console.log("MediaDataSource",a),O(a)},O=function(e){var t=document.getElementsByClassName("centeredVideo")[0];"undefined"!==typeof b&&null!=b&&(b.pause(),b.unload(),b.detachMediaElement(),b.destroy(),b=null),(b=$.a.createPlayer(e,{enableWorker:!1,lazyLoadMaxDuration:180,seekType:"range"})).on($.a.Events.ERROR,(function(e,t,a){console.log("errorType:",e),console.log("errorDetail:",t),console.log("errorInfo:",a),b&&(b.pause(),b.unload(),b.detachMediaElement(),b.destroy(),b=null,window.setTimeout(f,500))})),b.on($.a.Events.STATISTICS_INFO,(function(e){0!=h?h!=e.decodedFrames?h=e.decodedFrames:(console.log("decodedFrames:",e.decodedFrames),h=0,b&&(b.pause(),b.unload(),b.detachMediaElement(),b.destroy(),b=null,window.setTimeout(f,500))):h=e.decodedFrames})),b.attachMediaElement(t),b.load(),b.play()};return Object(D.jsx)("div",{children:Object(D.jsxs)(B.a,{fullScreen:!0,open:c,onClose:m,TransitionComponent:le,children:[Object(D.jsx)(U.a,{className:t.appBar,children:Object(D.jsxs)(q.a,{children:[Object(D.jsx)(T.a,{variant:"contained",onClick:f,children:"play"}),Object(D.jsxs)(z.a,{variant:"h6",className:t.title,children:["hasAudio",Object(D.jsx)(Y.a,{checked:u,id:"Audio",color:"primary",name:"hasAudio",onChange:function(e){j(e.target.checked)},inputProps:{"aria-label":"primary checkbox"}})]}),Object(D.jsx)(T.a,{autoFocus:!0,color:"inherit",onClick:m,children:Object(D.jsx)(_.a,{})})]})}),Object(D.jsx)("div",{className:t.videoContainer,children:Object(D.jsx)("div",{children:Object(D.jsx)("video",{name:"videoElement",className:"centeredVideo",controls:!0,allow:"autoPlay",width:"100%",children:"Your browser is too old which doesn't support HTML5 video."})})})]})})}var de=a(242),ue=a(240),je=a(221),be=Object(u.a)((function(e){return{appBar:{position:"relative"},title:{marginLeft:e.spacing(2),flex:1},formClass:{"& > *":{margin:e.spacing(1),width:"25ch"}},formDiv:{margin:"0 auto"}}})),he=o.a.forwardRef((function(e,t){return Object(D.jsx)(K.a,Object(H.a)({direction:"up",ref:t},e))}));function pe(e){var t=be(),a=o.a.useState(!1),r=Object(d.a)(a,2),c=r[0],i=r[1],l=o.a.useState(1===e.row.enabled),s=Object(d.a)(l,2),u=s[0],j=s[1],b=o.a.useState({id:e.row.id,code:e.row.code,rtmpAuthCode:e.row.rtmpAuthCode,playAuthCode:e.row.playAuthCode,onlineStatus:e.row.onlineStatus,enabled:e.row.enabled}),h=Object(d.a)(b,2),p=h[0],m=(h[1],o.a.useState(!1)),f=Object(d.a)(m,2),O=f[0],x=f[1],g=o.a.useState(""),v=Object(d.a)(g,2),w=v[0],y=v[1];Object(n.useImperativeHandle)(e.onRef,(function(){return{handleClickOpen:k}}));var k=function(){i(!0)},C=function(){i(!1)},S=function(t){oe(p).then((function(t){if(1===t.code)return i(!1),void(e.callBack&&e.callBack());y(t.msg),x(!0),window.setTimeout((function(){x(!1)}),5e3)}))},N=function(e){p[e.target.id]=e.target.value};return Object(D.jsx)("div",{children:Object(D.jsxs)(B.a,{fullScreen:!0,open:c,onClose:C,TransitionComponent:he,children:[Object(D.jsx)(U.a,{className:t.appBar,children:Object(D.jsxs)(q.a,{children:[Object(D.jsx)(T.a,{variant:"contained",onClick:S,children:"save"}),Object(D.jsx)("span",{children:"\xa0\xa0"}),"edit"===e.type?Object(D.jsx)(T.a,{variant:"contained",onClick:function(t){re(p).then((function(t){if(1===t.code)return i(!1),void(e.callBack&&e.callBack());y(t.msg),x(!0),window.setTimeout((function(){x(!1)}),5e3)}))},children:"delete"}):"",Object(D.jsx)(z.a,{variant:"h6",className:t.title}),Object(D.jsx)(T.a,{autoFocus:!0,color:"inherit",onClick:C,children:Object(D.jsx)(_.a,{})})]})}),O?Object(D.jsxs)(ue.a,{severity:"error",children:[Object(D.jsx)(je.a,{children:"Error"}),w," ",Object(D.jsx)("strong",{children:"check it out!"})]}):"",Object(D.jsx)("div",{className:t.formDiv,children:Object(D.jsxs)("form",{className:t.formClass,noValidate:!0,autoComplete:"off",onSubmit:S,children:["edit"===e.type?Object(D.jsx)("div",{children:Object(D.jsx)(de.a,{id:"id",label:"id",InputProps:{readOnly:!0},defaultValue:p.id})}):"",Object(D.jsx)("div",{children:Object(D.jsx)(de.a,{id:"code",label:"code",defaultValue:p.code,onChange:N})}),Object(D.jsx)("div",{children:Object(D.jsx)(de.a,{id:"rtmpAuthCode",label:"rtmpAuthCode",defaultValue:p.rtmpAuthCode,onChange:N})}),"edit"===e.type?"":Object(D.jsx)("div",{children:Object(D.jsx)(Y.a,{checked:u,id:"enabled",onChange:function(e){console.log(e.target.checked),p[e.target.id]=e.target.checked?1:0,"enabled"===e.target.id&&j(e.target.checked)},color:"primary",name:"enabled",inputProps:{"aria-label":"primary checkbox"}})})]})})]})})}var me=Object(u.a)((function(){return{root:{position:"relative"},dropdown:{position:"absolute",top:28,right:0,left:0,zIndex:999,border:"0px solid",padding:1,backgroundColor:"#f8f8f8"}}}));function fe(e){var t=me(),a=o.a.useState(!1),n=Object(d.a)(a,2),r=n[0],c=n[1],i=o.a.useState({title:"Success",content:"copy success !"}),l=Object(d.a)(i,2),s=l[0],u=(l[1],o.a.useState(e.row)),j=Object(d.a)(u,2),b=j[0],h=(j[1],o.a.createRef());var p=o.a.createRef();var m=o.a.createRef();return Object(D.jsx)(k.a,{onClickAway:function(){c(!1)},children:Object(D.jsxs)("div",{className:t.root,children:[Object(D.jsx)("button",{type:"button",onClick:function(){c((function(e){return!e}))},children:Object(D.jsx)(I.a,{})}),r?Object(D.jsx)("div",{className:t.dropdown,children:Object(D.jsxs)(C.a,{component:"nav","aria-label":"secondary mailbox folders",children:[Object(D.jsx)(S.a,{button:!0,children:Object(D.jsx)(N.a,{primary:"edit",onClick:function(){c(!1),p.current.handleClickOpen()}})}),Object(D.jsx)(S.a,{button:!0,onClick:function(){b.enabled=1===b.enabled?0:1,ce(b).then((function(t){if(1===t.code)return c(!1),void(e.callBack&&e.callBack())}))},children:1===e.row.enabled?Object(D.jsx)(N.a,{primary:"turn-off"}):Object(D.jsx)(N.a,{primary:"turn-on"})}),Object(D.jsx)(S.a,{button:!0,onClick:function(){c(!1),h.current.handleClickOpen()},children:Object(D.jsx)(N.a,{primary:"play"})}),Object(D.jsx)(S.a,{button:!0,onClick:function(){c(!1);var t=window.location.origin+window.location.pathname+"#/live?method=permanent&code="+e.row.code+"&authCode="+e.row.playAuthCode;R()(t),m.current.handleClickOpen()},children:Object(D.jsx)(N.a,{primary:"share"})})]})}):null,Object(D.jsx)(se,{row:e.row,onRef:h}),Object(D.jsx)(pe,{row:e.row,type:"edit",callBack:e.callBack,onRef:p}),Object(D.jsx)(V,{dialog:s,onRef:m})]})})}var Oe=[{id:"id",label:"id",minWidth:170},{id:"code",label:"code",minWidth:100},{id:"rtmpAuthCode",label:"rtmpAuthCode"},{id:"playAuthCode",label:"playAuthCode"},{id:"onlineStatus",label:"onlineStatus",format:function(e){return e&&1===e?Object(D.jsx)(y.a,{}):Object(D.jsx)(v.a,{})}},{id:"enabled",label:"enabled",format:function(e){return Object(D.jsx)(Y.a,{checked:1===e,id:"enabled",color:"primary",name:"enabled",inputProps:{"aria-label":"primary checkbox"}})}},{id:"action",label:"action",format:function(e,t,a){return Object(D.jsx)(fe,{row:t,callBack:a})}}],xe=Object(u.a)((function(e){return{root:{width:"100%",position:"relative"},container:{minHeight:400},dropdown:{position:"absolute",top:28,right:0,left:0,zIndex:0,border:"0px solid",padding:1,backgroundColor:"#f8f8f8"},appBar:{position:"relative"},title:{marginLeft:e.spacing(2),flex:1}}}));function ge(){var e,t,a=xe(),n=o.a.useState(0),r=Object(d.a)(n,2),c=r[0],i=r[1],l=o.a.useState(10),u=Object(d.a)(l,2),g=u[0],v=u[1],w=o.a.useState([]),y=Object(d.a)(w,2);e=y[0],t=y[1];var k=function(){ne().then((function(a){var n;1===a.code&&(e.splice(0),(n=e).push.apply(n,Object(s.a)(a.data.page)),t([]),t(e))}))},C=o.a.createRef();return o.a.useEffect(k,[]),Object(D.jsxs)(j.a,{className:a.root,children:[Object(D.jsx)(U.a,{className:a.appBar,children:Object(D.jsxs)(q.a,{children:[Object(D.jsx)(T.a,{variant:"contained",onClick:function(){C.current.handleClickOpen()},children:"ADD"}),Object(D.jsx)(z.a,{variant:"h6",className:a.title})]})}),Object(D.jsx)(pe,{row:{id:"",code:"",rtmpAuthCode:"",playAuthCode:"",onlineStatus:0,enabled:1},type:"add",callBack:k,onRef:C}),Object(D.jsx)(m.a,{className:a.container,children:Object(D.jsxs)(b.a,{stickyHeader:!0,"aria-label":"sticky table",children:[Object(D.jsx)(f.a,{children:Object(D.jsx)(x.a,{children:Oe.map((function(e){return Object(D.jsx)(p.a,{align:e.align,style:{minWidth:e.minWidth},children:e.label},e.id)}))})}),Object(D.jsx)(h.a,{children:e.slice(c*g,c*g+g).map((function(e){return Object(D.jsx)(x.a,{hover:!0,role:"checkbox",tabIndex:-1,children:Oe.map((function(t){var a=e[t.id];return Object(D.jsx)(p.a,{align:t.align,children:t.format?t.format(a,e,k):a},t.id)}))},e.code)}))})]})}),Object(D.jsx)(O.a,{rowsPerPageOptions:[10,25,100],component:"div",count:e.length,rowsPerPage:g,page:c,onChangePage:function(e,t){i(t)},onChangeRowsPerPage:function(e){v(+e.target.value),i(0)}})]})}var ve=a(246),we=a(236),ye=a(237),ke=a(244),Ce=a(234),Se=a(238),Ne=a(235),Ee=a(113),Ie=a.n(Ee);function Ae(){return Object(D.jsxs)(z.a,{variant:"body2",color:"textSecondary",align:"center",children:["Copyright \xa9 ",Object(D.jsx)(Ce.a,{color:"inherit",href:"https://material-ui.com/",children:"Your Website"})," ",(new Date).getFullYear(),"."]})}var Re=Object(u.a)((function(e){return{root:{height:"100vh"},image:{backgroundImage:"url(https://source.unsplash.com/random)",backgroundRepeat:"no-repeat",backgroundColor:"light"===e.palette.type?e.palette.grey[50]:e.palette.grey[900],backgroundSize:"cover",backgroundPosition:"center"},paper:{margin:e.spacing(8,4),display:"flex",flexDirection:"column",alignItems:"center"},avatar:{margin:e.spacing(1),backgroundColor:e.palette.secondary.main},form:{width:"100%",marginTop:e.spacing(1)},submit:{margin:e.spacing(3,0,2)}}}));function Te(){var e=Re(),t=o.a.useState({userName:"",password:""}),a=Object(d.a)(t,2),n=a[0],r=(a[1],o.a.useState(!1)),c=Object(d.a)(r,2),i=c[0],l=c[1],s=o.a.useState(""),u=Object(d.a)(s,2),b=u[0],h=u[1],p=function(e){n[e.target.id]=e.target.value};return o.a.useEffect((function(){"true"===window.localStorage.getItem("tokenExpired")&&localStorage.setItem("tokenExpired","false")}),[]),Object(D.jsxs)(Ne.a,{container:!0,component:"main",className:e.root,children:[Object(D.jsx)(we.a,{}),Object(D.jsx)(Ne.a,{item:!0,xs:!1,sm:4,md:7,className:e.image}),Object(D.jsx)(Ne.a,{item:!0,xs:12,sm:8,md:5,component:j.a,elevation:6,square:!0,children:Object(D.jsxs)("div",{className:e.paper,children:[Object(D.jsx)(ve.a,{className:e.avatar,children:Object(D.jsx)(Ie.a,{})}),Object(D.jsx)(z.a,{component:"h1",variant:"h5",children:"Sign in"}),i?Object(D.jsxs)(ue.a,{severity:"error",children:[Object(D.jsx)(je.a,{children:"Error"}),b," ",Object(D.jsx)("strong",{children:"check it out!"})]}):"",Object(D.jsxs)("form",{className:e.form,noValidate:!0,children:[Object(D.jsx)(de.a,{variant:"outlined",margin:"normal",required:!0,fullWidth:!0,id:"userName",label:"UserName",name:"userName",autoComplete:"userName",autoFocus:!0,onChange:p}),Object(D.jsx)(de.a,{variant:"outlined",margin:"normal",required:!0,fullWidth:!0,name:"password",label:"Password",type:"password",id:"password",autoComplete:"current-password",onChange:p}),Object(D.jsx)(ye.a,{control:Object(D.jsx)(ke.a,{value:"remember",color:"primary"}),label:"Remember me"}),Object(D.jsx)(T.a,{fullWidth:!0,variant:"contained",color:"primary",className:e.submit,onClick:function(e){ae(n).then((function(e){if(1===e.code){var t=window.localStorage;return t.setItem("token",e.data.token),t.setItem("tokenExpired","false"),void(window.location.hash="#/")}h(e.msg),l(!0)}))},children:"Sign In"}),Object(D.jsxs)(Ne.a,{container:!0,children:[Object(D.jsx)(Ne.a,{item:!0,xs:!0,children:Object(D.jsx)(Ce.a,{href:"#",variant:"body2",children:"Forgot password?"})}),Object(D.jsx)(Ne.a,{item:!0,children:Object(D.jsx)(Ce.a,{href:"#",variant:"body2",children:"Don't have an account? Sign Up"})})]}),Object(D.jsx)(Se.a,{mt:5,children:Object(D.jsx)(Ae,{})})]})]})})]})}var Be=a(7),Pe=a(8),Fe=a(21),Le=a(20),Me=function(e){Object(Fe.a)(a,e);var t=Object(Le.a)(a);function a(){return Object(Be.a)(this,a),t.apply(this,arguments)}return Object(Pe.a)(a,[{key:"render",value:function(){return Object(D.jsx)("h2",{children:"404"})}}]),a}(o.a.Component),De=Object(u.a)((function(e){return{appBar:{position:"relative",backgroundColor:"#eebbaa"},title:{marginLeft:e.spacing(2),flex:1},videoContainer:{width:"90%",margin:"0 auto"}}}));function We(e){var t=De(),a=o.a.useState(!0),n=Object(d.a)(a,2),r=n[0],c=n[1],i=null,l=0,s=function(e){var t=new RegExp("(^|&|\\?)"+e+"=([^&]*)(&|$)","i"),a=window.location.hash.substr(1).match(t);return null!=a?decodeURIComponent(a[2]):null},u=function(e){var t=s("method"),a=s("code"),n=s("authCode");if(t&&a&&n){var o={type:"flv"},c=te+"/live/"+t+"/"+a+"/"+n+".flv";o.url=c,o.hasAudio=r,o.isLive=!0,console.log("MediaDataSource",o),j(o)}},j=function(e){var t=document.getElementsByClassName("centeredVideo")[0];"undefined"!==typeof i&&null!=i&&(i.pause(),i.unload(),i.detachMediaElement(),i.destroy(),i=null),(i=$.a.createPlayer(e,{enableWorker:!1,lazyLoadMaxDuration:180,seekType:"range"})).on($.a.Events.ERROR,(function(e,t,a){console.log("errorType:",e),console.log("errorDetail:",t),console.log("errorInfo:",a),i&&(i.pause(),i.unload(),i.detachMediaElement(),i.destroy(),i=null,window.setTimeout(u,500))})),i.on($.a.Events.STATISTICS_INFO,(function(e){0!=l?l!=e.decodedFrames?l=e.decodedFrames:(console.log("decodedFrames:",e.decodedFrames),l=0,i&&(i.pause(),i.unload(),i.detachMediaElement(),i.destroy(),i=null,window.setTimeout(u,500))):l=e.decodedFrames})),i.attachMediaElement(t),i.load(),i.play()};return Object(D.jsx)("div",{children:Object(D.jsxs)("div",{children:[Object(D.jsx)(U.a,{className:t.appBar,children:Object(D.jsxs)(q.a,{children:[Object(D.jsx)(T.a,{variant:"contained",onClick:u,children:"play"}),Object(D.jsxs)(z.a,{variant:"h6",className:t.title,children:["hasAudio",Object(D.jsx)(Y.a,{checked:r,id:"Audio",color:"primary",name:"hasAudio",onChange:function(e){c(e.target.checked)},inputProps:{"aria-label":"primary checkbox"}})]})]})}),Object(D.jsx)("div",{className:t.videoContainer,children:Object(D.jsx)("div",{children:Object(D.jsx)("video",{name:"videoElement",className:"centeredVideo",controls:!0,allow:"autoPlay",width:"100%",children:"Your browser is too old which doesn't support HTML5 video."})})})]})})}function Ve(e){return Object(D.jsxs)(l.c,{children:[Object(D.jsx)(l.a,{exact:!0,path:"/",component:ge}),Object(D.jsx)(l.a,{exact:!0,path:"/live",component:We}),Object(D.jsx)(l.a,{path:"/login",component:Te}),Object(D.jsx)(l.a,{component:Me})]})}var He=Object(n.memo)(Ve);var Ue=function(){return Object(D.jsx)("div",{children:Object(D.jsx)(i.a,{children:Object(D.jsx)(He,{})})})};function qe(){return Object(D.jsx)(Ue,{})}c.a.render(Object(D.jsx)(qe,{}),document.querySelector("#app"))}},[[163,1,2]]]);
//# sourceMappingURL=main.407342d8.chunk.js.map