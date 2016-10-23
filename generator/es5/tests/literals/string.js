$module('string', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.fastbox('hello world', $g.____testlib.basictypes.String);
      $t.fastbox("hi world", $g.____testlib.basictypes.String);
      $t.fastbox('single quote with "quoted"', $g.____testlib.basictypes.String);
      $t.fastbox("double quote with 'quoted'", $g.____testlib.basictypes.String);
      $t.fastbox("escaped \" quote", $g.____testlib.basictypes.String);
      $t.fastbox('escaped \' quote', $g.____testlib.basictypes.String);
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
