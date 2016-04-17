$module('string', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.box('hello world', $g.____testlib.basictypes.String);
      $t.box("hi world", $g.____testlib.basictypes.String);
      $t.box('single quote with "quoted"', $g.____testlib.basictypes.String);
      $t.box("double quote with 'quoted'", $g.____testlib.basictypes.String);
      $t.box("escaped \" quote", $g.____testlib.basictypes.String);
      $t.box('escaped \' quote', $g.____testlib.basictypes.String);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
