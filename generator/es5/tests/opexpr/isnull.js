$module('isnull', function () {
  var $static = this;
  $static.DoSomething = function (a) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.box(a == null, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
});
