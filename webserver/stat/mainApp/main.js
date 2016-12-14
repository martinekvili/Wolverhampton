(function(mainApp) {
  document.addEventListener('DOMContentLoaded', function() {
    ng.platformBrowserDynamic
      .platformBrowserDynamic()
      .bootstrapModule(mainApp.AppModule);
  });
})(window.mainApp || (window.mainApp = {}));