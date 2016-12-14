(function(mainApp) {
  mainApp.AppModule =
    ng.core.NgModule({
      imports: [ ng.platformBrowser.BrowserModule, ng.router.RouterModule.forRoot([{ path: 'users', component: mainApp.UsersComponent },]), ng.http.HttpModule ],
      declarations: [ mainApp.AppComponent, mainApp.UsersComponent, mainApp.SidebarComponent, mainApp.TopNavComponent ],
      bootstrap: [ mainApp.AppComponent ]
    })
    .Class({
      constructor: function() {}
    });
})(window.mainApp || (window.mainApp = {}));