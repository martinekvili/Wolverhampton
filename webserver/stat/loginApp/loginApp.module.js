(function(loginApp) {
  loginApp.AppModule =
    ng.core.NgModule({
      imports: [ ng.platformBrowser.BrowserModule, ng.forms.FormsModule, ng.http.HttpModule ],
      declarations: [ loginApp.AppComponent ],
      bootstrap: [ loginApp.AppComponent ]
    })
    .Class({
      constructor: function() {}
    });
})(window.loginApp || (window.loginApp = {}));