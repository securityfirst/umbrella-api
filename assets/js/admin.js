$(function() {
    SetupMetisMenu();
    LoadCorrectSidebar();
    SetupEkkoLightBox();
    SetupSummerNote();
});

//Loads the correct sidebar on window load,
//collapses the sidebar on window resize.
// Sets the min-height of #page-wrapper to window size
function LoadCorrectSidebar() {
    $(window).bind("load resize", function() {
        topOffset = 50;
        width = (this.window.innerWidth > 0) ? this.window.innerWidth : this.screen.width;
        if (width < 768) {
            $('div.navbar-collapse').addClass('collapse');
            topOffset = 100; // 2-row-menu
        } else {
            $('div.navbar-collapse').removeClass('collapse');
        }

        height = ((this.window.innerHeight > 0) ? this.window.innerHeight : this.screen.height) - 1;
        height = height - topOffset;
        if (height < 1) height = 1;
        if (height > topOffset) {
            $("#page-wrapper").css("min-height", (height) + "px");
        }
    });

    var url = window.location;
    var element = $('ul.nav a').filter(function() {
        return this.href == url || url.href.indexOf(this.href) == 0;
    }).last();
    var element2 = element.addClass('active').parent().parent().addClass('in').parent();
    var element3 = element.parent().parent().parent().parent();
    if (element3.is('ul')) {
        element3.addClass('in');
    }
    if (element.is('li')) {
        element.addClass('active');
    }
};


function SetupEkkoLightBox() {
    $(document).delegate('*[data-toggle="lightbox"]', 'click', function(event) {
        event.preventDefault();
        $(this).ekkoLightbox();
    });
}

function SetupMetisMenu() {
    $('#side-menu').metisMenu();
}

function SetupSummerNote() {
    $('.summernote').summernote();
}
