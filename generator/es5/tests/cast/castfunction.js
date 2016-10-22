$module('castfunction', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.cast($g.castfunction.TEST, $g.____testlib.basictypes.Function($t.any), false);
      $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
});
