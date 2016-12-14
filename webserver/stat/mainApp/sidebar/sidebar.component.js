(function(mainApp) {
  mainApp.SidebarComponent =
    ng.core.Component({
      selector: 'sidebar-cmp',
      templateUrl: 'stat/mainApp/sidebar/template.html'
    })
    .Class({
      constructor: function() {}
    });
})(window.mainApp || (window.mainApp = {}));