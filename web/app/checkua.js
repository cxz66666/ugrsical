function check_ua(){
    return navigator.userAgent.toLocaleLowerCase().indexOf("micromessenger") > -1 || navigator.userAgent.indexOf("MQQBrowser") > -1 || navigator.userAgent.indexOf("QQTheme") > -1;
}