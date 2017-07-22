$module('full', function () {
  var $static = this;
  $static.TEST = function () {
    var lambda;
    lambda = function (firstParam, secondParam) {
      return secondParam;
    };
    return lambda($t.fastbox(123, $g.________testlib.basictypes.Integer), $t.fastbox(true, $g.________testlib.basictypes.Boolean));
  };
});
