$module('nativeboxing', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var r;
    var s;
    var sany;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            s = 'hello world';
            sany = s;
            r = $t.fastbox($t.cast(sany, $global.String, false), $g.____testlib.basictypes.String);
            $g.____testlib.basictypes.String.$equals(r, $t.fastbox('hello world', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
