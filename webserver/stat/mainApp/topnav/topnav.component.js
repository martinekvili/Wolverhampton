(function(mainApp) {
  mainApp.TopNavComponent =
    ng.core.Component({
      selector: 'top-nav',
      templateUrl: 'stat/mainApp/topnav/template.html'
    })
    .Class({
      constructor: function() {}
    });
})(window.mainApp || (window.mainApp = {}));