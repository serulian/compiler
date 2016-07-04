$module('mappingprops', function () {
  var $static = this;
  $static.SimpleFunction = function (props) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            props.$index($t.box('SomeProp', $g.____testlib.basictypes.String)).then(function ($result2) {
              return $promise.resolve($result2).then(function ($result1) {
                return $g.____testlib.basictypes.String.$equals($t.nullcompare($result1, $t.box('', $g.____testlib.basictypes.String)), $t.box('hello world', $g.____testlib.basictypes.String)).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
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
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).overObject(function () {
              var obj = {
              };
              obj["SomeProp"] = $t.box("hello world", $g.____testlib.basictypes.String);
              return obj;
            }()).then(function ($result1) {
              return $g.mappingprops.SimpleFunction($result1).then(function ($result0) {
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
