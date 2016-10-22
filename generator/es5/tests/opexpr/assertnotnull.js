$module('assertnotnull', function () {
  var $static = this;
  $static.TEST = function () {
    var someValue;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      someValue = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      $resolve($t.assertnotnull(someValue));
      return;
    };
    return $promise.new($continue);
  };
});
