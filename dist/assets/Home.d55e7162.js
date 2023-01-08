import{d as v,v as r,c as b,b as p,h as A,ac as te,r as ne,a7 as se,P as C,n as le,ad as ae,w as oe,y as re,x as ce,D as I,F as O,ae as ie,H as d,af as m,J as E,X as ue,G as de}from"./index.9a99308d.js";import{m as S,a as G,k as fe,u as ge,i as me,b as ve,l as he,c as _e,n as ye,o as Se,d as Ce,e as be,f as pe,p as ke,g as xe,q as we,r as Pe,h as $e,j as L,s as Ne}from"./VBtn.ddac2587.js";function Ve(e){for(;e;){if(Ee(e))return e;e=e.parentElement}return document.scrollingElement}function Ee(e){if(!e||e.nodeType!==Node.ELEMENT_NODE)return!1;const t=window.getComputedStyle(e);return t.overflowY==="scroll"||t.overflowY==="auto"&&e.scrollHeight>e.clientHeight}const Le=v({name:"VContainer",props:{fluid:{type:Boolean,default:!1},...S()},setup(e,t){let{slots:n}=t;return G(()=>r(e.tag,{class:["v-container",{"v-container--fluid":e.fluid}]},n)),{}}}),k=["sm","md","lg","xl","xxl"],H=(()=>k.reduce((e,t)=>(e[t]={type:[Boolean,String,Number],default:!1},e),{}))(),M=(()=>k.reduce((e,t)=>(e["offset"+b(t)]={type:[String,Number],default:null},e),{}))(),D=(()=>k.reduce((e,t)=>(e["order"+b(t)]={type:[String,Number],default:null},e),{}))(),j={col:Object.keys(H),offset:Object.keys(M),order:Object.keys(D)};function je(e,t,n){let l=e;if(!(n==null||n===!1)){if(t){const s=t.replace(e,"");l+=`-${s}`}return e==="col"&&(l="v-"+l),e==="col"&&(n===""||n===!0)||(l+=`-${n}`),l.toLowerCase()}}const Be=["auto","start","end","center","baseline","stretch"],B=v({name:"VCol",props:{cols:{type:[Boolean,String,Number],default:!1},...H,offset:{type:[String,Number],default:null},...M,order:{type:[String,Number],default:null},...D,alignSelf:{type:String,default:null,validator:e=>Be.includes(e)},...S()},setup(e,t){let{slots:n}=t;const l=p(()=>{const s=[];let c;for(c in j)j[c].forEach(i=>{const u=e[i],a=je(c,i,u);a&&s.push(a)});const f=s.some(i=>i.startsWith("v-col-"));return s.push({"v-col":!f||!e.cols,[`v-col-${e.cols}`]:e.cols,[`offset-${e.offset}`]:e.offset,[`order-${e.order}`]:e.order,[`align-self-${e.alignSelf}`]:e.alignSelf}),s});return()=>{var s;return A(e.tag,{class:l.value},(s=n.default)==null?void 0:s.call(n))}}}),Re=["sm","md","lg","xl","xxl"],x=["start","end","center"],z=["space-between","space-around","space-evenly"];function w(e,t){return Re.reduce((n,l)=>(n[e+b(l)]=t(),n),{})}const Te=[...x,"baseline","stretch"],F=e=>Te.includes(e),U=w("align",()=>({type:String,default:null,validator:F})),Ae=[...x,...z],Y=e=>Ae.includes(e),q=w("justify",()=>({type:String,default:null,validator:Y})),Ie=[...x,...z,"stretch"],J=e=>Ie.includes(e),W=w("alignContent",()=>({type:String,default:null,validator:J})),R={align:Object.keys(U),justify:Object.keys(q),alignContent:Object.keys(W)},Oe={align:"align",justify:"justify",alignContent:"align-content"};function Ge(e,t,n){let l=Oe[e];if(n!=null){if(t){const s=t.replace(e,"");l+=`-${s}`}return l+=`-${n}`,l.toLowerCase()}}const He=v({name:"VRow",props:{dense:Boolean,noGutters:Boolean,align:{type:String,default:null,validator:F},...U,justify:{type:String,default:null,validator:Y},...q,alignContent:{type:String,default:null,validator:J},...W,...S()},setup(e,t){let{slots:n}=t;const l=p(()=>{const s=[];let c;for(c in R)R[c].forEach(f=>{const i=e[f],u=Ge(c,f,i);u&&s.push(u)});return s.push({"v-row--no-gutters":e.noGutters,"v-row--dense":e.dense,[`align-${e.align}`]:e.align,[`justify-${e.justify}`]:e.justify,[`align-content-${e.alignContent}`]:e.alignContent}),s});return()=>{var s;return A(e.tag,{class:["v-row",l.value]},(s=n.default)==null?void 0:s.call(n))}}});function Me(e){return Math.floor(Math.abs(e))*Math.sign(e)}const De=v({name:"VParallax",props:{scale:{type:[Number,String],default:.5}},setup(e,t){let{slots:n}=t;const{intersectionRef:l,isIntersecting:s}=fe(),{resizeRef:c,contentRect:f}=ge(),{height:i}=te(),u=ne();se(()=>{var o;l.value=c.value=(o=u.value)==null?void 0:o.$el});let a;C(s,o=>{o?(a=Ve(l.value),a=a===document.scrollingElement?document:a,a.addEventListener("scroll",g,{passive:!0}),g()):a.removeEventListener("scroll",g)}),le(()=>{var o;(o=a)==null||o.removeEventListener("scroll",g)}),C(i,g),C(()=>{var o;return(o=f.value)==null?void 0:o.height},g);const h=p(()=>1-ae(+e.scale));let _=-1;function g(){!s.value||(cancelAnimationFrame(_),_=requestAnimationFrame(()=>{var N,V;var o;const P=((o=u.value)==null?void 0:o.$el).querySelector(".v-img__img");if(!P)return;const $=(N=a.clientHeight)!=null?N:document.documentElement.clientHeight,X=(V=a.scrollTop)!=null?V:window.scrollY,K=l.value.offsetTop,y=f.value.height,Q=K+(y-$)/2,Z=Me((X-Q)*h.value),ee=Math.max(1,(h.value*($-y)+y)/y);P.style.setProperty("transform",`translateY(${Z}px) scale(${ee})`)}))}return G(()=>r(me,{class:["v-parallax",{"v-parallax--active":s.value}],ref:u,cover:!0,onLoadstart:g,onLoad:g},n)),{}}});const T=v({name:"VSheet",props:{color:String,...ve(),...he(),..._e(),...ye(),...Se(),...Ce(),...S(),...oe()},setup(e,t){let{slots:n}=t;const{themeClasses:l}=re(e),{backgroundColorClasses:s,backgroundColorStyles:c}=be(ce(e,"color")),{borderClasses:f}=pe(e),{dimensionStyles:i}=ke(e),{elevationClasses:u}=xe(e),{locationStyles:a}=we(e),{positionClasses:h}=Pe(e),{roundedClasses:_}=$e(e);return()=>r(e.tag,{class:["v-sheet",l.value,s.value,f.value,u.value,h.value,_.value],style:[c.value,i.value,a.value]},n)}}),ze={class:"d-flex flex-column fill-height justify-center align-center text-white"},Fe=m("h1",{class:"text-h4 font-weight-thin mb-4"},"Hemanex CLI",-1),Ue=m("h4",{class:"subheading"},"Nexus Docker Registry CLI",-1),Ye={class:"mt-4"},qe=m("h1",null,"hello world",-1),Je=I({__name:"landing",setup(e){return(t,n)=>(O(),ie(ue,null,[r(De,{src:"https://assets.extrahop.com/images/blog/digital-wall.jpg",style:{top:"-65px"}},{default:d(()=>[m("div",ze,[Fe,Ue,m("div",Ye,[r(He,null,{default:d(()=>[r(B,null,{default:d(()=>[r(T,null,{default:d(()=>[r(L,{rounded:0,size:"large",color:"secondary",variant:"flat",elevation:"0"},{default:d(()=>[E(" Getting Started ")]),_:1})]),_:1})]),_:1}),r(B,null,{default:d(()=>[r(T,null,{default:d(()=>[r(L,{rounded:0,size:"large",color:"success",variant:"flat",elevation:"0"},{default:d(()=>[E(" Download Now ")]),_:1})]),_:1})]),_:1})]),_:1})])])]),_:1}),r(Le,{class:"fill-height"},{default:d(()=>[r(Ne,{class:"d-flex align-center text-center fill-height fill-width"},{default:d(()=>[qe]),_:1})]),_:1})],64))}}),Ke=I({__name:"Home",setup(e){return(t,n)=>(O(),de(Je))}});export{Ke as default};
