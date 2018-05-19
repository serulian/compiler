$module('tagged', function () {
  var $static = this;
  this.$struct('bad44f92', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        somefield: SomeField,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'somefield', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|fb1385bf<bad44f92>": true,
        "equals|4|fb1385bf<71258460>": true,
        "Stringify|2|fb1385bf<b2b53db7>": true,
        "Mapping|2|fb1385bf<204295f9<any>>": true,
        "Clone|2|fb1385bf<bad44f92>": true,
        "String|2|fb1385bf<b2b53db7>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
    var $result;
    var jsonString;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            s = $g.tagged.SomeStruct.new($t.fastbox(2, $g.________testlib.basictypes.Integer));
            jsonString = $t.fastbox('{"somefield":2}', $g.________testlib.basictypes.String);
            $promise.maybe(s.Stringify($g.________testlib.basictypes.JSON)()).then(function ($result0) {
              $result = $g.________testlib.basictypes.String.$equals($result0, jsonString);
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
  });
});
