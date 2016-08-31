$module('nochildren', function () {
  var $static = this;
  $static.SimpleFunction = function (props, children) {
    var $result;
    var $temp0;
    var $temp1;
    var found;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            found = $t.box(true, $g.____testlib.basictypes.Boolean);
            $current = 1;
            continue;

          case 1:
            $temp1 = children;
            $current = 2;
            continue;

          case 2:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $result;
            value = $temp0.First;
            if ($t.unbox($temp0.Second)) {
              $current = 4;
              continue;
            } else {
              $current = 5;
              continue;
            }
            break;

          case 4:
            found = $t.box(false, $g.____testlib.basictypes.Boolean);
            $current = 2;
            continue;

          case 5:
            $resolve(found);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty().then(function ($result1) {
              return $g.nochildren.SimpleFunction($result1, $generator.directempty()).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
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
