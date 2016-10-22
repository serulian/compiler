$module('escapedtemplatestr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.fastbox("hello 'world'! \"This is a\n\tlong quote!\"", $g.____testlib.basictypes.String);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
