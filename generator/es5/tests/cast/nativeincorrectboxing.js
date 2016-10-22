$module('nativeincorrectboxing', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var err;
    var result;
    var sany;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            sany = $t.fastbox('hello world', $g.____testlib.basictypes.String);
            $promise.new(function ($resolve) {
              $resolve($t.fastbox($t.cast(sany, $global.String, false), $g.____testlib.basictypes.String));
            }).then(function ($result0) {
              result = $result0;
              err = null;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function ($rejected) {
              err = $rejected;
              result = null;
              $current = 1;
              $continue($resolve, $reject);
              return;
            });
            return;

          case 1:
            $promise.resolve(result == null).then(function ($result0) {
              $result = $t.fastbox($result0 && (err != null), $g.____testlib.basictypes.Boolean);
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
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
