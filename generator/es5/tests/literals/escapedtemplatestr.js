$module('escapedtemplatestr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($continue) {
      $t.box("hello 'world'! \"This is a\n\tlong quote!\"", $g.____testlib.basictypes.String);
      $state.resolve();
    });
    return $promise.build($state);
  };
});
