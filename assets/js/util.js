//Get a pseudo UUID
uuid_counter = 0;
UUID = function() {
    uuid_counter += 1;
    return uuid_counter;
}

String.prototype.toCamelCase = function() {
  str = this.replace(/(\_[a-z])/g, function($1) {
      return $1.toUpperCase().replace('_', '');
  })

  return str
}

String.prototype.toClassCase = function() {
    var str = this.toCamelCase()
    return str.charAt(0).toUpperCase() + str.slice(1);
}

String.prototype.toUnderscore = function(){
	return this.replace(/([A-Z])/g, function($1){return "_"+$1.toLowerCase();});
};


String.prototype.lowercaseFirstLetter = function() {
 return this.charAt(0).toLowerCase() + this.slice(1);
}

$.fn.extend({
    animateCss: function (animationName, cb) {
        var animationEnd = 'webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend';
        $(this).addClass('animated ' + animationName).one(animationEnd, function() {
            $(this).removeClass('animated ' + animationName);
            if (cb !== null && cb !== undefined) {
              cb.apply(this);
            }
        });
    }
});

String.prototype.toProperCase = function () {
    var str = this.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});

    str = str.replace(new RegExp(/_/g), ' ')

    return str;
};

String.prototype.toTitleCase = function() {
  var i, j, str, lowers, uppers;
  var str = this.toProperCase();

  str = str.replace(/([^\W_]+[^\s-]*) */g, function(txt) {
    return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
  });

  // Certain minor words should be left lowercase unless 
  // they are the first or last words in the string
  lowers = ['A', 'An', 'The', 'And', 'But', 'Or', 'For', 'Nor', 'As', 'At', 
  'By', 'For', 'From', 'In', 'Into', 'Near', 'Of', 'On', 'Onto', 'To', 'With'];
  for (i = 0, j = lowers.length; i < j; i++)
    str = str.replace(new RegExp('\\s' + lowers[i] + '\\s', 'g'), 
      function(txt) {
        return txt.toLowerCase();
      });

  // Certain words such as initialisms or acronyms should be left uppercase
  uppers = ['Id', 'Tv'];
  for (i = 0, j = uppers.length; i < j; i++)
    str = str.replace(new RegExp('\\b' + uppers[i] + '\\b', 'g'), 
      uppers[i].toUpperCase());

  return str;
}
