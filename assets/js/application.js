require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
$(() => {
  // enabling bootstrap widgets
  $('[data-toggle="popover"]').popover();
  $('[data-toggle="tooltip"]').tooltip();

  // auto-close alerts
  window.setTimeout(function() {
    $(".alert:not('.alert-danger')").alert('close');
  }, 10000);

});
