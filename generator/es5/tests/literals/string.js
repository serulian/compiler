$module('string', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($continue) {
      $t.box('hello world', $g.____testlib.basictypes.String);
      $t.box("hi world", $g.____testlib.basictypes.String);
      $t.box('single quote with "quoted"', $g.____testlib.basictypes.String);
      $t.box("double quote with 'quoted'", $g.____testlib.basictypes.String);
      $t.box("escaped \" quote", $g.____testlib.basictypes.String);
      $t.box('escaped \' quote', $g.____testlib.basictypes.String);
      $state.resolve();
    });
    return $promise.build($state);
  };
});
