(function(mainApp) {
  mainApp.AppComponent =
    ng.core.Component({
      selector: 'main-app',
      templateUrl: 'stat/mainApp/template.html'
    })
    .Class({
      constructor: function() {}
    });
})(window.mainApp || (window.mainApp = {}));