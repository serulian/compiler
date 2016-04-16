$module('string', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.box('hello world', $g.____testlib.basictypes.String);
            $t.box("hi world", $g.____testlib.basictypes.String);
            $t.box('single quote with "quoted"', $g.____testlib.basictypes.String);
            $t.box("double quote with 'quoted'", $g.____testlib.basictypes.String);
            $t.box("escaped \" quote", $g.____testlib.basictypes.String);
            $t.box('escaped \' quote', $g.____testlib.basictypes.String);
            $state.current = -1;
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
