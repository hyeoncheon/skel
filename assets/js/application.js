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

  // navigation position highlighter
  var current_path = document.location.pathname.replace(/\/$/, "");
  $(".nav-item").removeClass("active");
  $(".nav-item").each(function(index) {
    if ($(this).attr('href') == current_path) {
      $(this).addClass("active");
      return false; // exit the loop
    }
  });
  $(".dropdown-item").removeClass("active");
  $(".dropdown-item").each(function(index) {
    if ($(this).attr('href') == current_path) {
      $(this).addClass("active");
      $(this).parent().parent().addClass("active");
      return false; // exit the loop
    }
  });

});
