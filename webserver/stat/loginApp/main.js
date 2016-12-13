(function(loginApp) {
  document.addEventListener('DOMContentLoaded', function() {
    ng.platformBrowserDynamic
      .platformBrowserDynamic()
      .bootstrapModule(loginApp.AppModule);
  });
})(window.loginApp || (window.loginApp = {}));