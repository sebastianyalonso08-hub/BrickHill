#include <stdint.h>
#define WINAPI __attribute__((stdcall))
typedef uint32_t DWORD; typedef uint16_t WORD;
typedef void* (WINAPI *PFN)(void);
typedef void* (WINAPI *PFN_LoadLibraryA)(const char*);
typedef void* (WINAPI *PFN_GetProcAddress)(void*, const char*);
typedef int (WINAPI *PFN_getsockopt)(void*,int,int,char*,int*);
typedef int (WINAPI *PFN_connect)(void*,const void*,int);

typedef struct LIST_ENTRY { struct LIST_ENTRY *Flink; struct LIST_ENTRY *Blink; } LIST_ENTRY;
typedef struct { WORD len; WORD max; uint16_t *buf; } UNICODE_STRING;

static void* peb(){void* p;__asm__ __volatile__("movl %%fs:0x30,%0":"=r"(p));return p;}
static int eqw(const uint16_t*w,const char*a){int i=0;for(;w[i]&&a[i];i++){uint16_t c=w[i];if(c>='A'&&c<='Z')c+=32;char d=a[i];if(d>='A'&&d<='Z')d+=32;if(c!=(unsigned char)d)return 0;}return w[i]==0&&a[i]==0;}
static void* mod(const char*n){uint8_t*p=(uint8_t*)peb();if(!p)return 0;uint8_t*l=*(uint8_t**)(p+0x0c);if(!l)return 0;LIST_ENTRY*h=(LIST_ENTRY*)(l+0x0c);LIST_ENTRY*e=h->Flink;for(int i=0;e&&e!=h&&i<128;i++,e=e->Flink){uint8_t*x=(uint8_t*)e;void*b=*(void**)(x+0x18);UNICODE_STRING*u=(UNICODE_STRING*)(x+0x2c);if(u->buf&&eqw(u->buf,n))return b;}return 0;}
static void* exp(void*m,const char*n){if(!m)return 0;uint8_t*b=(uint8_t*)m;uint32_t pe=*(uint32_t*)(b+0x3c);uint8_t*nt=b+pe;uint32_t dd=*(uint32_t*)(nt+0x78);if(!dd)return 0;uint8_t*ed=b+dd;uint32_t cnt=*(uint32_t*)(ed+0x18),np_rva=*(uint32_t*)(ed+0x20),fp_rva=*(uint32_t*)(ed+0x1c),op_rva=*(uint32_t*)(ed+0x24);uint32_t*np=(uint32_t*)(b+np_rva);uint16_t*op=(uint16_t*)(b+op_rva);uint32_t*fp=(uint32_t*)(b+fp_rva);for(uint32_t i=0;i<cnt;i++){const char*s=(const char*)(b+np[i]);const char*a=n;while(*s&&*a&&*s==*a){s++;a++;}if(!*s&&!*a)return b+fp[op[i]];}return 0;}
static PFN_GetProcAddress gpa;static PFN_LoadLibraryA lla;static void*ws(){static void*m; if(m)return m;void*k=mod("kernel32.dll");if(!k)k=mod("KERNEL32.DLL");gpa=(PFN_GetProcAddress)exp(k,"GetProcAddress");lla=(PFN_LoadLibraryA)exp(k,"LoadLibraryA");m=mod("ws2_32.dll");if(!m&&lla)m=lla("ws2_32.dll");return m;}
static void* fn(const char*n){void*m=ws();return m&&gpa?gpa(m,n):0;}

// Network-only shim: redirect TCP connects made by the unchanged legacy client to the local bridge.
__declspec(dllexport) int WINAPI connect(void*s,const void*addr,int len){
  PFN_connect real=(PFN_connect)fn("connect"); if(!real)return -1;
  if(len>=8 && ((const unsigned char*)addr)[0]==2){
    PFN_getsockopt gs=(PFN_getsockopt)fn("getsockopt"); int typ=0,tl=sizeof(typ);
    if(gs && gs(s,0xffff,0x1008,(char*)&typ,&tl)==0 && typ==1){
      unsigned char a[32]; for(int i=0;i<len;i++)a[i]=((const unsigned char*)addr)[i];
      a[2]=0x19; a[3]=0x6E; a[4]=127; a[5]=0; a[6]=0; a[7]=1;
      return real(s,a,len);
    }
  }
  return real(s,addr,len);
}
__declspec(dllexport) int WINAPI DllMain(void*h,unsigned r,void*p){return 1;}
